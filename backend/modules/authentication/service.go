package authentication

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"

	"checkpoint/db"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func SignInService(ctx iris.Context) {

}

func SignOutService(token string) {

}

func SignUpService(data SignUpData) (model.Account, error) {

	db := db.GetClient()
	hash, _ := hashPassword(data.Password)

	result := model.Account{}
	account := model.Account{
		ID:       uuid.New(),
		Username: data.Username,
		Password: hash,
	}

	insertStmt := Account.INSERT(Account.ID, Account.Username, Account.Password).MODEL(account).RETURNING(Account.AllColumns)

	err := insertStmt.Query(db, &result)

	return result, err
}
