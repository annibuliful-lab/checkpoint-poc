package projectRole

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/utils"
	"time"

	"github.com/google/uuid"
)

type CreateProjectRoleData struct {
	ProjectId     uuid.UUID `json:"projectId"`
	Title         string    `json:"title"`
	PermissionIds []string  `json:"permissionIds"`
}

type UpdateProjectRoleData struct {
	ID            uuid.UUID `json:"id"`
	ProjectId     uuid.UUID `json:"projectId"`
	Title         string    `json:"string"`
	PermissionIds []string  `json:"permissionIds"`
}

type GetProjectRoleByIdData struct {
	ProjectId uuid.UUID `json:"projectId"`
	ID        uuid.UUID `json:"id"`
}

type GetProjectRolesData struct {
	ProjectId  uuid.UUID `json:"projectId"`
	Search     string    `json:"search"`
	pagination utils.OffsetPagination
}

type DeleteProjectRoleData struct {
	ProjectId uuid.UUID `json:"projectId"`
	ID        uuid.UUID `json:"id"`
}

type ProjectRoleResponse struct {
	ID          uuid.UUID            `json:"id"`
	ProjectId   uuid.UUID            `json:"projectId"`
	Title       string               `json:"title"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   *time.Time           `json:"updatedAt"`
	Permissions []PermissionResponse `json:"permissions"`
}

type PermissionResponse struct {
	RoleID  uuid.UUID              `json:"roleId"`
	ID      uuid.UUID              `json:"id"`
	Subject string                 `json:"subject"`
	Action  model.PermissionAction `json:"action"`
}
