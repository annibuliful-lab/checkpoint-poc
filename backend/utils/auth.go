package utils

import (
	"strings"

	"github.com/kataras/iris/v12"
)

type AuthenticationHeader struct {
	Authorization string `header:"Authorization"`
	ProjectId     string `header:"x-project-id"`
	Token         string
}

func GetAuthenticationHeaders(ctx iris.Context) AuthenticationHeader {
	authorization := ctx.GetHeader("authorization")
	projectId := ctx.GetHeader("x-project-id")

	return AuthenticationHeader{
		Authorization: authorization,
		ProjectId:     projectId,
		Token:         GetAuthToken(authorization),
	}
}

func GetAuthToken(authorization string) string {
	return strings.Replace(authorization, "Bearer ", "", 1)
}
