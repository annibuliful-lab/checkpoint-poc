package auth

import "github.com/google/uuid"

type VerifyProjectAccountData struct {
	ID        uuid.UUID `json:"id"`
	AccountId uuid.UUID `json:"accountId"`
	Role      *string   `json:"role"`
}

type VerifyProjectRoleData struct {
	ID        uuid.UUID `json:"id"`
	AccountId uuid.UUID `json:"accountId"`
	ProjectId uuid.UUID `json:"projectId"`
}
