package auth

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/jwt"
	"checkpoint/utils"
	"context"
	"errors"
	"log"
	"time"

	"github.com/goccy/go-json"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func VerifyStationApiAuthentication(apiKey string) (*StationApiAuthenticationContext, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	redisClient := db.GetRedisClient()

	cacheResult, err := redisClient.Get(ctx, apiKey).Result()

	if err != nil {
		log.Println("veritification-station-api-cache-get-error", err.Error())
	}

	if cacheResult != "" {
		var cacheStation StationApiAuthenticationContext

		err = json.Unmarshal([]byte(cacheResult), &cacheStation)
		if err != nil {
			log.Println("json-cache-error", err.Error())
		}

		return &cacheStation, nil
	}

	getStationStmt := table.StationLocationConfiguration.
		SELECT(table.StationLocation.ID, table.StationLocation.ProjectId).
		FROM(
			table.StationLocation.
				INNER_JOIN(
					table.StationLocationConfiguration,
					table.StationLocationConfiguration.StationLocationId.EQ(table.StationLocation.ID),
				),
		).
		WHERE(
			table.StationLocationConfiguration.ApiKey.EQ(pg.String(apiKey)).
				AND(table.StationLocation.DeletedAt.IS_NULL()),
		)

	stationLocation := model.StationLocation{}
	err = getStationStmt.Query(dbClient, &stationLocation)

	if err != nil && db.HasNoRow(err) {
		return nil, errors.New("api key is invalid")
	}

	if err != nil {
		log.Println("station-api-access-error", err.Error())
		return nil, utils.InternalServerError
	}

	stationId := stationLocation.ID.String()
	projectId := stationLocation.ProjectId.String()

	jsonData, err := json.Marshal(StationApiAuthenticationContext{
		StationId: stationId,
		ProjectId: projectId,
		ApiKey:    apiKey,
	})

	if err != nil {
		log.Println("json-error: ", err.Error())
	}

	err = redisClient.Set(ctx, apiKey, jsonData, 15*time.Minute).Err()

	if err != nil {
		log.Println("redis-error: ", err.Error())
	}

	return &StationApiAuthenticationContext{
		StationId: stationId,
		ProjectId: projectId,
	}, nil
}

func VerifyAuthentication(headers AuthorizationContext) error {
	ctx := context.Background()
	dbClient := db.GetPrimaryClient()
	redisClient := db.GetRedisClient()
	payload, err := jwt.VerifyToken(headers.Token)

	if err != nil && err.Error() == utils.TokenExpire.Error() {

		err = redisClient.
			Del(ctx, headers.Token).
			Err()

		if err != nil {
			log.Println("redis-delete-cache-error", err.Error())
		}

		updateStmt := table.SessionToken.
			UPDATE(table.SessionToken.Revoke).
			MODEL(model.SessionToken{Revoke: true}).
			WHERE(table.SessionToken.Token.EQ(pg.String(headers.Token)))

		_, err := updateStmt.Exec(dbClient)
		if err != nil {
			log.Println("update-session-token-error", err.Error())
		}

		return utils.GraphqlError{
			Code:    "Token expired",
			Message: "Token expired",
		}
	}

	if err != nil {
		log.Println("jwt-error", err.Error())

		return utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	result, err := redisClient.Get(ctx, headers.Token).Result()

	if err != nil {
		log.Println("verification-authentication-cache-get-error", err.Error())
	}

	if result != "" {
		var cacheAccount utils.AuthorizationData

		err = json.Unmarshal([]byte(result), &cacheAccount)
		if err != nil {
			log.Println("json-cache-error", err.Error())
		}

		if cacheAccount.AccountId != payload.AccountId {
			log.Println("Account-mismatch")
			return utils.GraphqlError{
				Code: utils.ForbiddenOperation.Error(),
			}
		}

		if !cacheAccount.IsActive {
			return utils.GraphqlError{
				Message: utils.ContactOwner.Error(),
			}
		}

		return nil
	}

	var account struct {
		model.AccountConfiguration
	}

	selectAccountStmt := pg.
		SELECT(table.AccountConfiguration.AllColumns).
		FROM(table.AccountConfiguration).
		WHERE(table.AccountConfiguration.AccountId.EQ(pg.UUID(payload.AccountId))).
		LIMIT(1)

	err = selectAccountStmt.Query(dbClient, &account)

	if err != nil {
		log.Println("select-account-error", err.Error())

		return utils.GraphqlError{
			Message: utils.InternalServerError.Error(),
		}
	}

	if !account.IsActive {
		return utils.GraphqlError{
			Message: utils.ContactOwner.Error(),
		}
	}

	jsonData, err := json.Marshal(utils.AuthorizationData{
		AccountId: payload.AccountId,
		IsActive:  account.IsActive,
	})

	if err != nil {
		log.Println("Json-error", err.Error())
	}

	err = redisClient.Set(ctx, headers.Token, jsonData, 0).Err()
	if err != nil {
		log.Println("Cache-error", err.Error())
	}

	return nil
}
func VerifyProjectStation(ctx context.Context, headers AuthorizationContext) error {
	if headers.ProjectId == "" || headers.StationId == "" {
		return utils.GraphqlError{
			Message: "project id and station id are required",
		}
	}
	redisClient := db.GetRedisClient()
	dbClient := db.GetPrimaryClient()
	key := "projectId:" + headers.ProjectId + ",stationId:" + headers.StationId
	result, err := redisClient.Get(ctx, key).Result()
	if err == nil && result == "true" {
		return nil
	}

	selectStationStmt := table.StationLocation.
		SELECT(table.StationLocation.ID).
		WHERE(table.StationLocation.ProjectId.EQ(
			pg.UUID(uuid.MustParse(headers.ProjectId)),
		).AND(table.StationLocation.DeletedAt.IS_NULL())).
		LIMIT(1)

	_, err = selectStationStmt.Exec(dbClient)

	if err != nil && db.HasNoRow(err) {
		return utils.GraphqlError{
			Message: "station id is not permitted to access",
		}
	}

	err = redisClient.Set(ctx, key, "true", 15*time.Minute).Err()
	if err != nil {
		log.Println("cache-project-station-error", err.Error())
	}

	return nil
}

func VerifyAuthorization(ctx context.Context, headers AuthorizationContext, permissionData utils.AuthorizationPermissionData) error {

	if headers.Token == "" || headers.ProjectId == "" {

		return utils.GraphqlError{
			Message: utils.InvalidToken.Error(),
		}
	}

	payload, _ := jwt.VerifyToken(headers.Token)
	key := "accountId:" + payload.AccountId.String() + "," + "projectId:" + headers.ProjectId

	// Check if the authorization data is present in the cache
	result, err := db.GetRedisClient().Get(ctx, key).Result()
	if err == nil && result != "" {
		cacheError := handleCachedAuthorization(result, permissionData)
		if cacheError != nil {
			return utils.GraphqlError{
				Message: cacheError.Error(),
			}
		}

		return nil
	}

	projectId, err := uuid.Parse(headers.ProjectId)
	if err != nil {

		return utils.GraphqlError{
			Message: "Project id is required",
		}
	}

	projectAccount, err := getProjectAccount(payload.AccountId, projectId)
	if err != nil {
		log.Println("get-project-account-error", err.Error())

		return utils.GraphqlError{
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	accountPermissions, err := getAccountPermissions(projectAccount.ProjectRole.ID)
	if err != nil {
		log.Println("get-account-permission-error", err.Error())
		return utils.GraphqlError{
			Message: utils.InternalServerError.Error(),
		}
	}

	if !hasPermission(accountPermissions, permissionData) {

		return utils.GraphqlError{
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	cacheAuthorization(ctx, key, payload.AccountId, projectId, accountPermissions)
	return nil
}

func handleCachedAuthorization(result string, permissionData utils.AuthorizationPermissionData) error {
	var data utils.AuthorizationWithPermissionsData
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		log.Println("Cache-error", err.Error())
		return utils.GraphqlError{
			Message: err.Error(),
		}
	}

	_, match := lo.Find(data.Permissions, func(el utils.AuthorizationPermissionData) bool {
		return el.PermissionAction == permissionData.PermissionAction && el.PermissionSubject == permissionData.PermissionSubject
	})

	if !match {

		return utils.GraphqlError{
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	return nil
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
			table.ProjectAccount.AccountId,
			table.ProjectRole.ID,
			table.ProjectRole.Title,
		).
		FROM(
			table.ProjectAccount.
				INNER_JOIN(table.ProjectRole, table.ProjectRole.ID.EQ(table.ProjectAccount.RoleId)),
		).
		WHERE(
			table.ProjectAccount.AccountId.EQ(pg.UUID(accountID)).
				AND(table.ProjectAccount.ProjectId.EQ(pg.UUID(projectID))))

	err := selectProjectAccountStmt.Query(dbClient, &projectAccount)
	return projectAccount, err
}

func getAccountPermissions(roleID uuid.UUID) ([]struct{ model.Permission }, error) {
	dbClient := db.GetPrimaryClient()
	var accountPermissions []struct{ model.Permission }

	selectProjectAccountPermissionsStmt := pg.
		SELECT(table.Permission.Action, table.Permission.Subject, table.Permission.ID).
		FROM(
			table.ProjectRolePermission.
				INNER_JOIN(table.Permission, table.Permission.ID.EQ(table.ProjectRolePermission.PermissionId)),
		).WHERE(table.ProjectRolePermission.RoleId.EQ(pg.UUID(roleID)))

	err := selectProjectAccountPermissionsStmt.Query(dbClient, &accountPermissions)

	return accountPermissions, err
}

func hasPermission(accountPermissions []struct{ model.Permission }, permissionData utils.AuthorizationPermissionData) bool {
	_, match := lo.Find(accountPermissions, func(el struct{ model.Permission }) bool {
		return el.Action == permissionData.PermissionAction && el.Subject == permissionData.PermissionSubject
	})
	return match
}

func cacheAuthorization(ctx context.Context, key string, accountID uuid.UUID, projectID uuid.UUID, accountPermissions []struct{ model.Permission }) {
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

	err = db.GetRedisClient().Set(ctx, key, jsonData, 15*time.Minute).Err()
	if err != nil {
		log.Println("Cache-error", err.Error())
	}
}

func GetStationAuthorizationContext(ctx context.Context) StationAuthorizationContext {

	return StationAuthorizationContext{
		StationId: ctx.Value("stationId").(string),
		ApiKey:    ctx.Value("apiKey").(string),
		DeviceId:  ctx.Value("deviceId").(string),
		ProjectId: ctx.Value("projectId").(string),
	}
}

func GetAuthorizationContext(ctx context.Context) AuthorizationContext {
	auth := AuthorizationContext{
		Token:     ctx.Value("token").(string),
		ProjectId: ctx.Value("projectId").(string),
		AccountId: ctx.Value("accountId").(string),
	}

	if ctx.Value("stationId") != nil {
		auth.StationId = ctx.Value("stationId").(string)
	}

	return auth
}
