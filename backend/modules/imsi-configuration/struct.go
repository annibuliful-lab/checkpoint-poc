package imsiconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/utils"
	"time"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type Imsiconfiguration struct {
	ID                graphql.ID                 `json:"id"`
	ProjectId         graphql.ID                 `json:"projectId"`
	Imsi              string                     `json:"imsi"`
	CreatedBy         string                     `json:"createdBy"`
	UpdatedBy         *string                    `json:"updatedBy"`
	CreatedAt         time.Time                  `json:"createdAt"`
	UpdatedAt         *time.Time                 `json:"updatedAt"`
	Label             model.DevicePermittedLabel `json:"label"`
	Priority          model.BlacklistPriority    `json:"priority"`
	StationLocationId uuid.UUID                  `json:"stationLocationId"`
	Mcc               string                     `json:"mcc"`
	Mnc               string                     `json:"mnc"`
}

type CreateImeiConfigurationInput struct {
	StationLocationId graphql.ID `json:"stationLocationId"`
	Imsi              string     `json:"imsi"`
	PermittedLabel    string     `json:"permittedLabel"`
	Priority          string     `json:"priority"`
	Mcc               string     `json:"mcc"`
	Mnc               string     `json:"mnc"`
	Tags              []string   `jcon:"tags"`
}

type CreateImsiConfigurationData struct {
	ProjectId         uuid.UUID `json:"projectId"`
	Imsi              string    `json:"imsi"`
	CreatedBy         string    `json:"createdBy"`
	PermittedLabel    string    `json:"label"`
	Priority          string    `json:"priority"`
	StationLocationId uuid.UUID `json:"stationLocationId"`
	Tags              []string  `json:"tags"`
}

type UpdateImsiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	Imsi      string    `json:"imsi"`
	ProjectId uuid.UUID `json:"projectId"`
	UpdatedBy string    `json:"updatedBy"`
	Label     string    `json:"label"`
	Priority  string    `json:"priority"`
	Tags      []string  `json:"tags"`
}

type GetImsiConfigurationsData struct {
	Search     string    `json:"search"`
	ProjectId  uuid.UUID `json:"projectId"`
	Label      string    `json:"label"`
	Tags       []string  `json:"tags"`
	Mnc        string    `json:"mnc"`
	Mcc        string    `json:"mcc"`
	Pagination utils.OffsetPagination
}

type GetImsiConfigurationByIdData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
}

type DeleteImsiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	DeletedBy string    `json:"updatedBy"`
}
