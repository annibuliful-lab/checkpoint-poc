package middleware

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/jwt"
	"checkpoint/utils"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	lo "github.com/samber/lo"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/kataras/iris/v12"
)

func AuthMiddleware() iris.Handler {
	dbClient := db.GetPrimaryClient()
	redisClient := db.GetRedisClient()

	return func(ctx iris.Context) {
		headers := utils.GetAuthenticationHeaders(ctx)

		if headers.Token == "" {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": utils.InvalidToken.Error(),
			})
			return
		}

		payload, err := jwt.VerifyToken(headers.Token)

		if err != nil && err.Error() == utils.TokenExpire.Error() {

			err = redisClient.
				Del(ctx, headers.Token).
				Err()

			if err != nil {
				log.Println("redis-delete-cache-error", err.Error())
			}

			updateStmt := SessionToken.
				UPDATE(SessionToken.Revoke).
				MODEL(model.SessionToken{Revoke: true}).
				WHERE(SessionToken.Token.EQ(pg.String(headers.Token)))

			_, err := updateStmt.Exec(dbClient)
			if err != nil {
				log.Println("update-session-token-error", err.Error())
			}

			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": "Token expired",
			})
			return
		}

		if err != nil {
			log.Println("jwt-error", err.Error())
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": err.Error(),
			})
			return
		}

		result, err := redisClient.Get(ctx, utils.GetAuthToken(headers.Authorization)).Result()

		if err != nil {
			log.Println("Cache-get-error", err.Error())
		}

		if result != "" {
			var cacheAccount utils.AuthorizationData

			err = json.Unmarshal([]byte(result), &cacheAccount)
			if err != nil {
				log.Println("Json-cache-error", err.Error())
			}

			if cacheAccount.AccountId != payload.AccountId {
				log.Println("Account-mismatch")
				ctx.StatusCode(iris.StatusForbidden)
				ctx.JSON(iris.Map{
					"message": utils.ForbiddenOperation.Error(),
				})
				return
			}

			if !cacheAccount.IsActive {
				ctx.StatusCode(iris.StatusForbidden)
				ctx.JSON(iris.Map{
					"message": utils.ContactOwner,
				})
				return
			}

			ctx.Next()
			return
		}

		var account struct {
			model.AccountConfiguration
		}

		selectAccountStmt := pg.
			SELECT(AccountConfiguration.AllColumns).
			FROM(AccountConfiguration).
			WHERE(AccountConfiguration.AccountId.EQ(pg.UUID(payload.AccountId))).
			LIMIT(1)

		err = selectAccountStmt.Query(dbClient, &account)

		if err != nil {
			log.Println(err.Error())
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"message": utils.InternalServerError,
			})
			return
		}

		if !account.IsActive {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": utils.ContactOwner,
			})
			return
		}

		jsonData, err := json.Marshal(utils.AuthorizationData{
			AccountId: payload.AccountId,
			IsActive:  account.IsActive,
		})
		if err != nil {
			log.Println("Json-error", err.Error())
		}

		err = redisClient.Set(ctx, utils.GetAuthToken(headers.Authorization), jsonData, 0).Err()
		if err != nil {
			log.Println("Cache-error", err.Error())
		}

		ctx.Next()
	}
}

func VerifyAuthorizationMiddleware(permissionData utils.AuthorizationPermissionData) iris.Handler {
	return func(ctx iris.Context) {
		headers := utils.GetAuthenticationHeaders(ctx)

		if headers.Token == "" || headers.ProjectId == "" {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": utils.InvalidToken.Error(),
			})
			return
		}

		payload, _ := jwt.VerifyToken(headers.Token)
		key := "accountId:" + payload.AccountId.String() + "," + "projectId:" + headers.ProjectId

		// Check if the authorization data is present in the cache
		result, err := db.GetRedisClient().Get(ctx, key).Result()
		if err == nil && result != "" {
			handleCachedAuthorization(ctx, result, permissionData)
			return
		}

		projectId, err := uuid.Parse(headers.ProjectId)
		if err != nil {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": "project id is required",
			})
			return
		}

		projectAccount, err := getProjectAccount(payload.AccountId, projectId)
		if err != nil {
			log.Println("get-project-account-error", err.Error())
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": utils.ForbiddenOperation.Error(),
			})
			return
		}

		accountPermissions, err := getAccountPermissions(projectAccount.ProjectRole.ID)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"message": utils.InternalServerError.Error(),
			})
			return
		}

		if !hasPermission(accountPermissions, permissionData) {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": utils.ForbiddenOperation.Error(),
			})
			return
		}

		// Cache the authorization data
		cacheAuthorization(ctx, key, payload.AccountId, projectId, accountPermissions)

		ctx.Next()
	}
}

func handleCachedAuthorization(ctx iris.Context, result string, permissionData utils.AuthorizationPermissionData) {
	var data utils.AuthorizationWithPermissionsData
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		log.Println("Cache-error", err.Error())
		return
	}

	_, match := lo.Find(data.Permissions, func(el utils.AuthorizationPermissionData) bool {
		return el.PermissionAction == permissionData.PermissionAction && el.PermissionSubject == permissionData.PermissionSubject
	})

	if !match {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.ForbiddenOperation.Error(),
		})
		return
	}

	ctx.Next()
}

func getProjectAccount(accountID uuid.UUID, projectID uuid.UUID) (struct {
	model.ProjectAccount
	model.ProjectRole
}, error) {
	dbClient := db.GetPrimaryClient()
	var projectAccount struct {
		model.ProjectAccount
		model.ProjectRole
	}

	selectProjectAccountStmt := pg.
		SELECT(
			ProjectAccount.AccountId,
			ProjectRole.ID,
			ProjectRole.Title,
		).
		FROM(
			ProjectAccount.
				INNER_JOIN(ProjectRole, ProjectRole.ID.EQ(ProjectAccount.RoleId)),
		).
		WHERE(
			ProjectAccount.AccountId.EQ(pg.UUID(accountID)).
				AND(ProjectAccount.ProjectId.EQ(pg.UUID(projectID))))

	err := selectProjectAccountStmt.Query(dbClient, &projectAccount)
	return projectAccount, err
}

func getAccountPermissions(roleID uuid.UUID) ([]struct{ model.Permission }, error) {
	dbClient := db.GetPrimaryClient()
	var accountPermissions []struct{ model.Permission }

	selectProjectAccountPermissionsStmt := pg.
		SELECT(Permission.Action, Permission.Subject, Permission.ID).
		FROM(
			ProjectRolePermission.
				INNER_JOIN(Permission, Permission.ID.EQ(ProjectRolePermission.PermissionId)),
		).WHERE(ProjectRolePermission.RoleId.EQ(pg.UUID(roleID)))

	err := selectProjectAccountPermissionsStmt.Query(dbClient, &accountPermissions)

	return accountPermissions, err
}

func hasPermission(accountPermissions []struct{ model.Permission }, permissionData utils.AuthorizationPermissionData) bool {
	_, match := lo.Find(accountPermissions, func(el struct{ model.Permission }) bool {
		return el.Action == permissionData.PermissionAction && el.Subject == permissionData.PermissionSubject
	})
	return match
}

func cacheAuthorization(ctx iris.Context, key string, accountID uuid.UUID, projectID uuid.UUID, accountPermissions []struct{ model.Permission }) {
	data := utils.AuthorizationWithPermissionsData{
		AccountId: accountID,
		ProjectId: projectID,
		Permissions: lo.Map(accountPermissions, func(item struct{ model.Permission }, index int) utils.AuthorizationPermissionData {
			return utils.AuthorizationPermissionData{
				PermissionSubject: item.Subject,
				PermissionAction:  item.Action,
			}
		}),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Json-error", err.Error())
		return
	}

	err = db.GetRedisClient().Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		log.Println("Cache-error", err.Error())
	}
}
