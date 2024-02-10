package authentication

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	jwt "checkpoint/jwt"
	utils "checkpoint/utils"
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/graph-gophers/graphql-go"
	"github.com/thanhpk/randstr"

	"github.com/google/uuid"
)

type AuthenticationService struct{}

func (AuthenticationService) SignIn(data SignInData) (*Authentication, string, error) {
	dbClient := db.GetPrimaryClient()
	redisClient := db.GetRedisClient()

	var account struct {
		model.Account
		model.AccountConfiguration
	}

	selectStmt := pg.SELECT(Account.AllColumns, AccountConfiguration.AllColumns).
		FROM(
			Account.
				INNER_JOIN(
					AccountConfiguration,
					Account.ID.EQ(AccountConfiguration.AccountId),
				),
		).
		WHERE(Account.Username.EQ(pg.String(data.Username))).
		LIMIT(1)

	err := selectStmt.Query(dbClient, &account)

	if err != nil && db.HasNoRow(err) {
		log.Println(err.Error())
		return nil, "username or password is incorrect", errors.New("username or password is incorrect")
	}

	if err != nil {
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	match, err := comparePasswordAndHash(data.Password, account.Password)

	if err != nil {
		log.Println(err.Error())
		return nil, "username or password is incorrect", errors.New("username or password is incorrect")

	}
	if !match {
		return nil, "username or password is incorrect", errors.New("username or password is incorrect")
	}

	if !account.AccountConfiguration.IsActive {
		return nil, "please contact project owner", errors.New("please contact project owner")
	}

	token, err := jwt.SignToken(jwt.SignedTokenParams{
		AccountId: account.ID.String(),
		Nounce:    randstr.Hex(16),
	})

	if err != nil {
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.SignTokenFailed
	}

	refreshToken, err := jwt.SignRefreshToken(jwt.SignedTokenParams{
		AccountId: account.ID.String(),
		Nounce:    randstr.Hex(16),
	})

	if err != nil {
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.SignTokenFailed
	}

	tx, err := dbClient.Begin()

	if err != nil {
		log.Println("init-db-tx", err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	insertSessionTokenStmt := SessionToken.
		INSERT(SessionToken.AccountId, SessionToken.Token, SessionToken.IsRefreshToken).
		MODEL(model.SessionToken{
			AccountId:      account.ID,
			Token:          token,
			IsRefreshToken: false,
		}).
		MODEL(model.SessionToken{
			AccountId:      account.ID,
			Token:          refreshToken,
			IsRefreshToken: true,
		})

	ctx := context.Background()

	_, err = insertSessionTokenStmt.ExecContext(ctx, tx)

	if err != nil {
		log.Println("insert-error: ", err.Error())
		tx.Rollback()
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	jsonData, err := json.Marshal(utils.AuthorizationData{
		AccountId: account.AccountId,
		IsActive:  account.IsActive,
	})

	if err != nil {
		log.Println("json-error: ", err.Error())
		tx.Rollback()
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	err = redisClient.Set(ctx, token, jsonData, time.Minute*15).Err()

	if err != nil {
		log.Println("redis-error: ", err.Error())
		tx.Rollback()
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	err = tx.Commit()
	if err != nil {
		log.Println("insert-session-token-error", err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	return &Authentication{
		UserId:       graphql.ID(account.ID.String()),
		Token:        token,
		RefreshToken: refreshToken,
	}, "Signin", nil
}

func (AuthenticationService) SignOut(token string) error {
	dbClient := db.GetPrimaryClient()
	redisClient := db.GetRedisClient()
	ctx := context.Background()

	updateStmt := SessionToken.
		UPDATE(SessionToken.Revoke).
		MODEL(model.SessionToken{
			Revoke: true,
		}).
		WHERE(SessionToken.Token.EQ(pg.String(token)))

	_, err := updateStmt.Exec(dbClient)
	if err != nil {
		log.Println(err.Error())
		return utils.InternalServerError
	}

	err = redisClient.
		Del(ctx, token).
		Err()

	if err != nil {
		log.Println(err.Error())
		return utils.InternalServerError
	}

	return nil
}

func (AuthenticationService) SignUp(data SignUpData) (model.Account, string, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	hash, _ := hashPassword(data.Password)
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println("init-db-tx", err.Error())
		return model.Account{}, utils.InternalServerError.Error(), utils.InternalServerError
	}

	accountResult := model.Account{}
	insertAccountStmt := Account.
		INSERT(Account.ID, Account.Username, Account.Password, Account.CreatedBy).
		MODEL(model.Account{
			ID:        uuid.New(),
			Username:  data.Username,
			Password:  hash,
			CreatedBy: "System",
		}).
		RETURNING(Account.AllColumns)

	err = insertAccountStmt.QueryContext(ctx, tx, &accountResult)

	if err != nil {
		tx.Rollback()

		if strings.Contains(err.Error(), "duplicate") {
			return model.Account{}, utils.DataConflict.Error(), errors.New("username is already exists")
		}

		return model.Account{}, utils.InternalServerError.Error(), utils.InternalServerError
	}

	accountConfiguration := model.AccountConfiguration{}
	insertAccountConfigurationStmt := AccountConfiguration.
		INSERT(AccountConfiguration.AccountId, AccountConfiguration.IsActive).
		MODEL(model.AccountConfiguration{
			AccountId: accountResult.ID,
			IsActive:  true,
		}).
		RETURNING(AccountConfiguration.AllColumns)

	err = insertAccountConfigurationStmt.QueryContext(ctx, tx, &accountConfiguration)
	if err != nil {
		tx.Rollback()
		return model.Account{}, utils.InternalServerError.Error(), utils.InternalServerError
	}

	err = tx.Commit()
	if err != nil {
		log.Println("signup-error", err.Error())
		return model.Account{}, utils.InternalServerError.Error(), utils.InternalServerError
	}

	return accountResult, "Signup", nil
}

func (AuthenticationService) GetAuthenticationTokenByRefreshToken(data RefreshTokenData) (*SigninResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	redisClient := db.GetRedisClient()

	selectSessionTokenStmt := SessionToken.
		UPDATE(SessionToken.Revoke).
		MODEL(model.SessionToken{Revoke: true}).
		WHERE(pg.
			AND(SessionToken.IsRefreshToken.EQ(pg.Bool(true)),
				SessionToken.Revoke.NOT_EQ(pg.Bool(true)),
				SessionToken.Token.EQ(pg.String(data.RefreshToken))),
		).
		RETURNING(SessionToken.AllColumns)

	sessionToken := model.SessionToken{}

	err := selectSessionTokenStmt.Query(dbClient, &sessionToken)

	if err != nil && db.HasNoRow(err) {
		return nil, 403, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println(err.Error())
		return nil, 500, utils.InternalServerError
	}

	accountId := sessionToken.AccountId.String()

	selectAccountConfigurationStmt := pg.
		SELECT(AccountConfiguration.IsActive, AccountConfiguration.AccountId).
		FROM(AccountConfiguration).
		WHERE(AccountConfiguration.AccountId.EQ(pg.UUID(uuid.MustParse(accountId))))

	var accountConfiguration struct {
		model.AccountConfiguration
	}

	err = selectAccountConfigurationStmt.Query(dbClient, &accountConfiguration)

	if err != nil {
		log.Println("select-account-error: ", err.Error())
		return nil, 500, utils.InternalServerError
	}

	token, err := jwt.SignToken(jwt.SignedTokenParams{
		AccountId: accountId,
		Nounce:    randstr.Hex(16),
	})

	if err != nil {
		log.Println(err.Error())
		return nil, 500, utils.SignTokenFailed
	}

	refreshToken, err := jwt.SignRefreshToken(jwt.SignedTokenParams{
		AccountId: accountId,
		Nounce:    randstr.Hex(16),
	})

	if err != nil {
		log.Println(err.Error())
		return nil, 500, utils.SignTokenFailed
	}

	tx, err := dbClient.Begin()

	if err != nil {
		log.Println("init-db-tx", err.Error())
		return nil, 500, utils.InternalServerError
	}

	insertSessionTokenStmt := SessionToken.
		INSERT(SessionToken.AccountId, SessionToken.Token, SessionToken.IsRefreshToken).
		MODEL(model.SessionToken{
			AccountId:      sessionToken.AccountId,
			Token:          token,
			IsRefreshToken: false,
		}).
		MODEL(model.SessionToken{
			AccountId:      sessionToken.AccountId,
			Token:          refreshToken,
			IsRefreshToken: true,
		})
	ctx := context.Background()

	_, err = insertSessionTokenStmt.ExecContext(ctx, tx)
	if err != nil {
		log.Println("insert-token: ", err.Error())
		tx.Rollback()
		return nil, 500, utils.InternalServerError
	}

	jsonData, err := json.Marshal(utils.AuthorizationData{
		AccountId: accountConfiguration.AccountId,
		IsActive:  accountConfiguration.IsActive,
	})

	if err != nil {
		log.Println("json-error: ", err.Error())
		tx.Rollback()
		return nil, 500, utils.InternalServerError
	}

	err = redisClient.Set(ctx, token, jsonData, time.Minute*15).Err()

	if err != nil {
		log.Println("redis-error: ", err.Error())
		tx.Rollback()
		return nil, 500, utils.InternalServerError
	}

	err = tx.Commit()
	if err != nil {
		log.Println("signin-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	return &SigninResponse{
		UserId:       sessionToken.AccountId.String(),
		Token:        token,
		RefreshToken: refreshToken,
	}, 200, nil
}
