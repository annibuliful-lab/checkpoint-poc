package projectRole

import (
	"checkpoint/utils"
	"time"

	"github.com/google/uuid"
)

type CreateProjectRole struct {
	Title string `json:"string"`
}

type UpdateProjectRole struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"string"`
}

type GetProjectRoleById struct {
	ID uuid.UUID `json:"id"`
}

type DeleteProjectRole struct {
	ID uuid.UUID `json:"id"`
}

type GetProjectRole struct {
	Search *string `json:"search"`
	utils.OffsetPagination
}

type ProjectRoleResponse struct {
	ID        uuid.UUID  `json:"id"`
	ProjectId uuid.UUID  `json:"projectId"`
	Title     string     `json:"string"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
