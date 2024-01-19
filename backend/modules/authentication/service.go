package authentication

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	jwt "checkpoint/jwt"
	utils "checkpoint/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	pg "github.com/go-jet/jet/v2/postgres"

	"github.com/google/uuid"
)

func SignInService(data SignInData) (SigninResponse, error) {
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
		return SigninResponse{}, errors.New(utils.InternalServerError)
	}

	if !account.AccountConfiguration.IsActive {
		return SigninResponse{}, errors.New("please contact project owner")
	}

	match, err := comparePasswordAndHash(data.Password, account.Password)

	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, errors.New("username or password is incorrect")

	}
	if !match {
		return SigninResponse{}, errors.New("username or password is incorrect")
	}

	token, err := jwt.SignToken(jwt.SignedTokenParams{
		AccountId: account.ID.String(),
	})

	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, errors.New(utils.SignTokenFailed)
	}

	refreshToken, err := jwt.SignRefreshToken(jwt.SignedTokenParams{
		AccountId: account.ID.String(),
	})

	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, errors.New(utils.SignTokenFailed)
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
		return SigninResponse{}, errors.New(utils.InternalServerError)
	}

	err = redisClient.
		Set(ctx, token, fmt.Sprintf("%+v", model.Account{
			ID: account.ID,
		}), 0).
		Err()

	if err != nil {
		log.Println(err.Error())
		return SigninResponse{}, errors.New(utils.InternalServerError)
	}

	return SigninResponse{
		UserId:       account.ID.String(),
		Token:        token,
		RefreshToken: refreshToken,
	}, err
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
		return errors.New(utils.InternalServerError)
	}

	err = redisClient.
		Del(ctx, token).
		Err()

	if err != nil {
		log.Println(err.Error())
		return errors.New(utils.InternalServerError)
	}

	return nil
}

func SignUpService(data SignUpData) (model.Account, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	hash, _ := hashPassword(data.Password)
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println(err.Error())
		return model.Account{}, errors.New(utils.InternalServerError)
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
			return model.Account{}, errors.New("username is already exists")
		}

		return model.Account{}, errors.New(utils.InternalServerError)
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
		return model.Account{}, errors.New(utils.InternalServerError)
	}

	tx.Commit()

	return accountResult, err
}
