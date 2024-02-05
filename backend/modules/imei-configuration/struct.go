package imeiconfiguration

import (
	"checkpoint/utils"
	"time"

	"github.com/google/uuid"
)

type CreateImeiConfigurationData struct {
	ProjectId         uuid.UUID `json:"projectId"`
	Imei              string    `json:"imei"`
	Priority          string    `json:"priority"`
	StationLocationId uuid.UUID `json:"stationLocationId"`
	CreatedBy         string    `json:"createdBy"`
	Label             string    `json:"label"`
	Tags              []string  `json:"tags"`
}

type UpdateImeiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	Imei      string    `json:"imei"`
	Priority  string    `json:"priority"`
	UpdatedBy string    `json:"createdBy"`
	Label     string    `json:"label"`
	Tags      []string  `json:"tags"`
}

type DeleteImeiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	DeletedBy string    `json:"deletedBy"`
}

type GetImeiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
}

type GetImeiConfigurationsData struct {
	ProjectId  uuid.UUID `json:"projectId"`
	Search     string    `json:"search"`
	Label      string    `json:"label"`
	Tags       []string  `json:"tags"`
	Pagination utils.OffsetPagination
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
