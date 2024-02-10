package imeiconfiguration

import (
	"checkpoint/utils"
	"time"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type CreateImeiConfigurationInput struct {
	Imei              string     `json:"imei"`
	Priority          string     `json:"priority"`
	StationLocationId graphql.ID `json:"stationLocationId"`
	PermittedLabel    string     `json:"permittedLabel"`
	Tags              *[]string  `json:"tags"`
}

type CreateImeiConfigurationData struct {
	ProjectId         uuid.UUID `json:"projectId"`
	Imei              string    `json:"imei"`
	Priority          string    `json:"priority"`
	StationLocationId uuid.UUID `json:"stationLocationId"`
	CreatedBy         string    `json:"createdBy"`
	Label             string    `json:"label"`
	Tags              *[]string `json:"tags"`
}

type UpdateImeiConfigurationInput struct {
	ID             graphql.ID `json:"id"`
	Imei           string     `json:"imei"`
	Priority       string     `json:"priority"`
	PermittedLabel string     `json:"permittedLabel"`
	Tags           *[]string  `json:"tags"`
}

type UpdateImeiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	Imei      string    `json:"imei"`
	Priority  string    `json:"priority"`
	UpdatedBy string    `json:"createdBy"`
	Label     string    `json:"label"`
	Tags      *[]string `json:"tags"`
}

type DeleteImeiConfigurationInput struct {
	ID graphql.ID `json:"id"`
}

type DeleteImeiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	DeletedBy string    `json:"deletedBy"`
}

type GetImeiConfigurationInput struct {
	ID graphql.ID `json:"id"`
}

type GetImeiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
}

type GetImeiConfigurationsInput struct {
	StationLocationId graphql.ID `json:"stationLocationId"`
	Search            *string    `json:"search"`
	Label             *string    `json:"label"`
	Tags              *[]string  `json:"tags"`
	Limit             float64    `json:"limit"`
	Skip              float64    `json:"skip"`
}

type GetImeiConfigurationsData struct {
	StationLocationId uuid.UUID `json:"stationLocationId"`
	ProjectId         uuid.UUID `json:"projectId"`
	Search            *string   `json:"search"`
	Label             *string   `json:"label"`
	Tags              *[]string `json:"tags"`
	Pagination        utils.OffsetPagination
}

type ImeiConfiguration struct {
	ID                graphql.ID        `json:"id"`
	ProjectId         graphql.ID        `json:"projectId"`
	Imei              string            `json:"imei"`
	CreatedBy         string            `json:"createdBy"`
	UpdatedBy         *graphql.NullID   `json:"updatedBy"`
	CreatedAt         graphql.Time      `json:"createdAt"`
	UpdatedAt         *graphql.NullTime `json:"updatedAt"`
	Priority          string            `json:"priority"`
	StationLocationId graphql.ID        `json:"stationLocationId"`
	PermittedLabel    string            `json:"permittedLabel"`
}

type ImeiConfigurationResponse struct {
	ID                uuid.UUID  `json:"id"`
	ProjectId         uuid.UUID  `json:"projectId"`
	Imei              string     `json:"imei"`
	CreatedBy         string     `json:"createdBy"`
	UpdatedBy         *string    `json:"updatedBy"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	Priority          string     `json:"priority"`
	StationLocationId uuid.UUID  `json:"stationLocationId"`
	Label             string     `json:"label"`
	Tags              []string   `json:"tags"`
}
