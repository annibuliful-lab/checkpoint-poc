package authentication

import (
	"checkpoint/jwt"
	"checkpoint/utils"

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

	account, code, err := SignUpService(data)

	if err != nil {
		ctx.StatusCode(code)
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

	response, code, err := SignInService(SignInData{
		Username: data.Username,
		Password: data.Password,
	})

	if err != nil {
		ctx.StatusCode(code)
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

func RefreshTokenController(ctx iris.Context) {
	refreshToken := ctx.GetHeader("authorization")

	if refreshToken == "" {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.InvalidToken,
		})
		return
	}

	token := utils.GetAuthToken(refreshToken)
	_, jwtError := jwt.VerifyToken(token)

	if jwtError != nil {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.InvalidToken,
		})
		return
	}

	response, code, err := GetAuthenticationTokenByRefreshToken(RefreshTokenData{
		RefreshToken: token,
	})

	if err != nil {
		ctx.StatusCode(code)
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

	token := utils.GetAuthToken(authorizationToken)

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
