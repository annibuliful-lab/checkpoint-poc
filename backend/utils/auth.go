package utils

import (
	"checkpoint/.gen/checkpoint/public/model"
	"strings"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type AuthenticationHeader struct {
	Authorization string `header:"Authorization"`
	ProjectId     string `header:"x-project-id"`
	Token         string
}

type AuthorizationData struct {
	AccountId uuid.UUID `json:"accountId"`
	IsActive  bool      `json:"isActive"`
}

type AuthorizationWithPermissionsData struct {
	AccountId   uuid.UUID                     `json:"accountId"`
	ProjectId   uuid.UUID                     `json:"projectId"`
	Permissions []AuthorizationPermissionData `json:"permissions"`
}

type AuthorizationPermissionData struct {
	PermissionSubject string
	PermissionAction  model.PermissionAction
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
