package middleware

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/jwt"
	"checkpoint/utils"
	"encoding/json"
	"log"
	"strings"

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
			updateStmt := SessionToken.
				UPDATE(SessionToken.Revoke).
				MODEL(model.SessionToken{Revoke: true}).
				WHERE(SessionToken.Token.EQ(pg.String(headers.Token)))

			_, err := updateStmt.Exec(dbClient)
			if err != nil {
				log.Println(err.Error())
			}

		}

		if err != nil {
			log.Println(err.Error())
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

func VerifyAuthorizationMiddleware(param utils.AuthorizationPermissionData) iris.Handler {
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

		payload, _ := jwt.VerifyToken(headers.Token)

		key := "accountId:" + payload.AccountId.String() + "," + "projectId" + headers.ProjectId

		result, err := redisClient.Get(ctx, key).Result()

		if err != nil {
			log.Println("Cache-get-error", err.Error())
		}

		if result != "" {
			var data utils.AuthorizationWithPermissionsData
			err = json.Unmarshal([]byte(result), &data)
			if err != nil {
				log.Println("Cache-error", err.Error())
			}

			_, match := lo.Find(data.Permissions, func(el utils.AuthorizationPermissionData) bool {
				return el.PermissionAction == param.PermissionAction && el.PermissionSubject == param.PermissionSubject
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
				ProjectAccount.AccountId.EQ(pg.UUID(payload.AccountId)).
					AND(ProjectAccount.ProjectId.EQ(pg.UUID(uuid.MustParse(headers.ProjectId)))))

		err = selectProjectAccountStmt.Query(dbClient, &projectAccount)

		if err != nil && strings.Contains(err.Error(), "no rows") {
			log.Println(err.Error())
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": utils.ForbiddenOperation.Error(),
			})
			return
		}

		var accountPermissions []struct {
			model.Permission
		}

		selectProjectAccountPermisionsStmt := pg.
			SELECT(Permission.Action, Permission.Subject).
			FROM(
				ProjectRolePermission.
					INNER_JOIN(Permission, Permission.ID.EQ(ProjectRolePermission.PermissionId)),
			).WHERE(ProjectRolePermission.RoleId.EQ(pg.UUID(projectAccount.ProjectRole.ID)))

		err = selectProjectAccountPermisionsStmt.Query(dbClient, &accountPermissions)

		if err != nil {
			log.Println(err.Error())
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"message": utils.InternalServerError.Error(),
			})
		}

		if len(accountPermissions) == 0 {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": utils.ForbiddenOperation.Error(),
			})
			return
		}

		_, match := lo.Find(accountPermissions, func(el struct{ model.Permission }) bool {
			return el.Action == param.PermissionAction && el.Subject == param.PermissionSubject
		})

		if !match {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": utils.ForbiddenOperation.Error(),
			})
			return
		}

		data := utils.AuthorizationWithPermissionsData{
			AccountId: payload.AccountId,
			ProjectId: uuid.MustParse(headers.ProjectId),
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
		}

		err = redisClient.Set(ctx, key, jsonData, 0).Err()
		if err != nil {
			log.Println("Cache-error", err.Error())
		}

		ctx.Next()
	}
}
