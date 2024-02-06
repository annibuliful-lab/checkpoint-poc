package projectRole

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/utils"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type ProjectRole struct {
	Id        graphql.ID   `json:"id"`
	ProjectId graphql.ID   `json:"projectId"`
	Title     string       `json:"title"`
	CreatedAt graphql.Time `json:"createdAt"`
}

type ProjectRolePermission struct {
	Id      graphql.ID `json:"id"`
	RoleId  graphql.ID `json:"roleId"`
	Subject string     `json:"subject"`
	Action  string     `json:"action"`
}

type CreateProjectRoleInput struct {
	Title         string       `json:"title"`
	PermissionIds []graphql.ID `json:"permissionIds"`
}

type UpdateProjectRoleInput struct {
	Id            graphql.ID   `json:"id"`
	Title         string       `json:"title"`
	PermissionIds []graphql.ID `json:"permissionIds"`
}

type CreateProjectRoleData struct {
	ProjectId     uuid.UUID    `json:"projectId"`
	Title         string       `json:"title"`
	PermissionIds []graphql.ID `json:"permissionIds"`
}

type UpdateProjectRoleData struct {
	ID            uuid.UUID    `json:"id"`
	ProjectId     uuid.UUID    `json:"projectId"`
	Title         string       `json:"string"`
	PermissionIds []graphql.ID `json:"permissionIds"`
}

type GetProjectRoleByIdData struct {
	ProjectId uuid.UUID `json:"projectId"`
	ID        uuid.UUID `json:"id"`
}

type GetProjectRolesInput struct {
	Search *string `json:"search"`
	Limit  int32   `json:"limit"`
	Skip   int32   `json:"skip"`
}

type GetProjectRolesData struct {
	ProjectId  uuid.UUID `json:"projectId"`
	Search     *string   `json:"search"`
	pagination utils.OffsetPagination
}

type DeleteProjectRoleData struct {
	ProjectId uuid.UUID  `json:"projectId"`
	ID        graphql.ID `json:"id"`
	AccountId uuid.UUID  `json:"accountId"`
}

type PermissionResponse struct {
	RoleID  uuid.UUID              `json:"roleId"`
	ID      uuid.UUID              `json:"id"`
	Subject string                 `json:"subject"`
	Action  model.PermissionAction `json:"action"`
}
