package authentication

import "github.com/kataras/iris/v12"

func SignUpController(ctx iris.Context) {
	response := SignUp(SignUpData{
		Username: "awdasd",
		Password: "wasdasd",
	})

	ctx.JSON(response)
}

func SignInController(ctx iris.Context) {

}

func SignOutController(ctx iris.Context) {

}
