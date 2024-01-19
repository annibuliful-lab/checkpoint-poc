package authentication

import (
	"checkpoint/jwt"
	"checkpoint/utils"
	"strings"

	"github.com/kataras/iris/v12"
)

func SignUpController(ctx iris.Context) {
	var data SignUpData
	err := ctx.ReadJSON(&data)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

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
		ctx.StatusCode(iris.StatusBadRequest)
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
	authorizationToken := ctx.GetHeader("authorization")

	if authorizationToken == "" {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.InvalidToken,
		})
		return
	}

	token := strings.Replace(authorizationToken, "Bearer", "", 1)

	_, jwtError := jwt.VerifyToken(token)

	if jwtError != nil {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.InvalidToken,
		})
		return
	}

	err := SignOutService(token)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"message": err.Error(),
		})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"message": "success",
	})
}
