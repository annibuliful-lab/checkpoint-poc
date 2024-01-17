package authentication

import (
	"github.com/kataras/iris/v12"
)

func SignUpController(ctx iris.Context) {
	var data SignUpData
	err := ctx.ReadJSON(&data)

	account, err := SignUpService(SignUpData{
		Username: data.Username,
		Password: data.Password,
	})

	if err != nil {
		ctx.StatusCode(500)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.StatusCode(201)
	ctx.JSON(iris.Map{
		"message": "created",
		"data": iris.Map{
			"id":        account.ID,
			"createdAt": account.CreatedAt,
			"updatedAt": account.UpdatedAt,
		},
	})
}

func SignInController(ctx iris.Context) {
	var data SignInData
	err := ctx.ReadJSON(&data)
	if err != nil {
		ctx.StatusCode(500)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	response, err := SignInService(SignInData{
		Username: data.Username,
		Password: data.Password,
	})

	if err != nil {
		ctx.StatusCode(500)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"message": "signed",
		"data":    response,
	})

}

func SignOutController(ctx iris.Context) {

}
