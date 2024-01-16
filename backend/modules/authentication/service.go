package authentication

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"

	"checkpoint/db"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type Authentication struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
}

type SignUpData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignIn(ctx iris.Context) {

}

func SignOut(token string) {

}

func SignUp(data SignUpData) model.Account {

	db := db.GetClient()
	account := model.Account{
		ID:       uuid.New(),
		Username: data.Username,
		Password: data.Password,
	}

	insertStmt := Account.INSERT(Account.ID, Account.Username, Account.Password).MODEL(account).RETURNING(Account.ID, Account.Username, Account.Password)

	result := model.Account{}

	err := insertStmt.Query(db, &result)

	if err == nil {
		panic(err.Error())
	}

	return result
}
