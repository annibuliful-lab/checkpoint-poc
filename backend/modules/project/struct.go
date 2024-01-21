package project

import (
	"time"

	"github.com/google/uuid"
)

type ProjectResponse struct {
	Title     string     `json:"title"`
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type CreateProjectData struct {
	Title     string `json:"title"`
	AccountId string `json:"accountId"`
}

type UpdateProjectData struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type GetProjectData struct {
	ID uuid.UUID `json:"id"`
}

type DeleteProjectData struct {
	ID        uuid.UUID `json:"id"`
	AccountId *string   `json:"accountId"`
}

type VerifyProjectAccountData struct {
	ID        uuid.UUID `json:"id"`
	AccountId uuid.UUID `json:"accountId"`
	Role      *string   `json:"role"`
}
