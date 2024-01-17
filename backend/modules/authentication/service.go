package authentication

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"context"
	"errors"
	"fmt"
	"strings"

	"checkpoint/db"
	. "checkpoint/jwt"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func SignInService(data SignInData) (SigninResponse, error) {
	db := db.GetClient()

	var account struct {
		model.Account
		model.AccountConfiguration
	}

	selectStmt := SELECT(Account.AllColumns, AccountConfiguration.AllColumns).
		FROM(
			Account.
				INNER_JOIN(
					AccountConfiguration,
					Account.ID.EQ(AccountConfiguration.AccountId),
				),
		).
		WHERE(Account.Username.EQ(String(data.Username))).
		LIMIT(1)

	err := selectStmt.Query(db, &account)
	if err != nil {

		fmt.Println(err.Error())
		return SigninResponse{}, errors.New("internal server error")
	}

	if account.AccountConfiguration.IsActive == false {
		return SigninResponse{}, errors.New("please contact project owner")
	}

	match, err := comparePasswordAndHash(data.Password, account.Password)
	if match == false {
		return SigninResponse{}, errors.New("username or password is incorrect")
	}

	token, err := SignToken(SignedTokenParams{
		AccountId: account.ID.String(),
	})

	if err != nil {
		return SigninResponse{}, errors.New("Sign token is failed")
	}

	refreshToken, err := SignRefreshToken(SignedTokenParams{
		AccountId: account.ID.String(),
	})

	if err != nil {
		return SigninResponse{}, errors.New("Sign token is failed")
	}

	return SigninResponse{
		UserId:       account.ID.String(),
		Token:        token,
		RefreshToken: refreshToken,
	}, err
}

func SignOutService(token string) {

}

func SignUpService(data SignUpData) (model.Account, error) {
	ctx := context.Background()
	db := db.GetClient()
	hash, _ := hashPassword(data.Password)
	tx, err := db.Begin()
	account := model.Account{
		ID:       uuid.New(),
		Username: data.Username,
		Password: hash,
	}

	accountResult := model.Account{}
	insertAccountStmt := Account.
		INSERT(Account.ID, Account.Username, Account.Password).
		MODEL(account).
		RETURNING(Account.AllColumns)
	err = insertAccountStmt.QueryContext(ctx, tx, &accountResult)

	if err != nil {
		tx.Rollback()

		if strings.Contains(err.Error(), "duplicate") {
			return model.Account{}, errors.New("username is already exists")
		}

		return model.Account{}, errors.New("internal server error")
	}

	accountConfiguration := model.AccountConfiguration{
		AccountId: accountResult.ID,
		IsActive:  true,
	}
	accountConfigurationResult := model.AccountConfiguration{}
	insertAccountConfigurationStmt := AccountConfiguration.
		INSERT(AccountConfiguration.AccountId, AccountConfiguration.IsActive).
		MODEL(accountConfiguration).
		RETURNING(AccountConfiguration.AllColumns)
	err = insertAccountConfigurationStmt.QueryContext(ctx, tx, &accountConfigurationResult)
	if err != nil {
		tx.Rollback()
		return model.Account{}, errors.New("internal server error")
	}

	tx.Commit()

	return accountResult, err
}
