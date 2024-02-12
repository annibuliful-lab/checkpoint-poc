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

type AuthenticationHeader struct {
	Authorization string `json:"authorization"`
	ProjectId     string `json:"projectId"`
	Token         string `json:"token"`
}

type AuthorizationContext struct {
	Token     string `json:"token"`
	ProjectId string `json:"projectId"`
	AccountId string `json:"accountId"`
}
