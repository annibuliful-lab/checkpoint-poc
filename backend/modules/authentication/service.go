package authentication

import "github.com/kataras/iris/v12"

type Authentication struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
}

func SignIn(ctx iris.Context) {

}

func SignOut(ctx iris.Context) {

}
