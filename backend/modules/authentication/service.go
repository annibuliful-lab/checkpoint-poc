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

	pg "github.com/go-jet/jet/v2/postgres"

	"github.com/google/uuid"
)

func SignInService(data SignInData) (SigninResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	redisClient := db.GetRedisClient()
	ctx := context.Background()

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
	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, 500, utils.InternalServerError
	}

	match, err := comparePasswordAndHash(data.Password, account.Password)

	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, 401, errors.New("username or password is incorrect")

	}
	if !match {
		return SigninResponse{}, 401, errors.New("username or password is incorrect")
	}

	if !account.AccountConfiguration.IsActive {
		return SigninResponse{}, 403, errors.New("please contact project owner")
	}

	token, err := jwt.SignToken(jwt.SignedTokenParams{
		AccountId: account.ID.String(),
	})

	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, 500, utils.SignTokenFailed
	}

	refreshToken, err := jwt.SignRefreshToken(jwt.SignedTokenParams{
		AccountId: account.ID.String(),
	})

	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, 500, utils.SignTokenFailed
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

	_, err = insertSessionTokenStmt.Exec(dbClient)
	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, 500, utils.InternalServerError
	}

	jsonData, err := json.Marshal(utils.AuthorizationData{
		AccountId: account.AccountId,
		IsActive:  account.IsActive,
	})

	if err != nil {
		log.Println("Json-error", err.Error())
		return SigninResponse{}, 500, utils.InternalServerError
	}

	err = redisClient.Set(ctx, token, jsonData, 0).Err()

	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, 500, utils.InternalServerError
	}

	return SigninResponse{
		UserId:       account.ID.String(),
		Token:        token,
		RefreshToken: refreshToken,
	}, 200, err
}

func SignOutService(token string) error {
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

func SignUpService(data SignUpData) (model.Account, int, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	hash, _ := hashPassword(data.Password)
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println(err.Error())
		return model.Account{}, 500, utils.InternalServerError
	}

	accountResult := model.Account{}
	insertAccountStmt := Account.
		INSERT(Account.ID, Account.Username, Account.Password).
		MODEL(model.Account{
			ID:       uuid.New(),
			Username: data.Username,
			Password: hash,
		}).
		RETURNING(Account.AllColumns)
	err = insertAccountStmt.QueryContext(ctx, tx, &accountResult)

	if err != nil {
		tx.Rollback()

		if strings.Contains(err.Error(), "duplicate") {
			return model.Account{}, 409, errors.New("username is already exists")
		}

		return model.Account{}, 500, utils.InternalServerError
	}

	accountConfigurationResult := model.AccountConfiguration{}
	insertAccountConfigurationStmt := AccountConfiguration.
		INSERT(AccountConfiguration.AccountId, AccountConfiguration.IsActive).
		MODEL(model.AccountConfiguration{
			AccountId: accountResult.ID,
			IsActive:  true,
		}).
		RETURNING(AccountConfiguration.AllColumns)

	err = insertAccountConfigurationStmt.QueryContext(ctx, tx, &accountConfigurationResult)
	if err != nil {
		tx.Rollback()
		return model.Account{}, 500, utils.InternalServerError
	}

	tx.Commit()

	return accountResult, 200, err
}
