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
	ApiKey        string `json:"apiKey"`
	DeviceId      string `json:"deviceId"`
	StationId     string `json:"stationId"`
}

type AuthorizationContext struct {
	Token     string `json:"token"`
	ProjectId string `json:"projectId"`
	AccountId string `json:"accountId"`
	StationId string `json:"stationId"`
}

type StationAuthorizationContext struct {
	StationId string `json:"stationId"`
	ApiKey    string `json:"apiKey"`
	DeviceId  string `json:"deviceId"`
	ProjectId string `json:"projectId"`
}

type StationAuthContext struct {
	ApiKey   string `json:"apiKey"`
	DeviceId string `json:"deviceId"`
}

type StationApiAuthenticationContext struct {
	StationId string `json:"stationId"`
	ApiKey    string `json:"apiKey"`
	ProjectId string `json:"projectId"`
}
