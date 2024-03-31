package imsiconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/gql/enum"
	"checkpoint/utils"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type ImsiConfiguration struct {
	ID                graphql.ID                `json:"id"`
	ProjectId         graphql.ID                `json:"projectId"`
	Imsi              string                    `json:"imsi"`
	CreatedBy         graphql.ID                `json:"createdBy"`
	UpdatedBy         *graphql.NullID           `json:"updatedBy"`
	CreatedAt         graphql.Time              `json:"createdAt"`
	UpdatedAt         *graphql.NullTime         `json:"updatedAt"`
	PermittedLabel    enum.DevicePermittedLabel `json:"label"`
	BlacklistPriority enum.BlacklistPriority    `json:"priority"`
	StationLocationId graphql.ID                `json:"stationLocationId"`
	Mcc               string                    `json:"mcc"`
	Mnc               string                    `json:"mnc"`
}

type UpsertImsiConfigurationInput struct {
	Imsi              string                     `json:"imsi"`
	PermittedLabel    model.DevicePermittedLabel `json:"permittedLabel"`
	BlacklistPriority model.BlacklistPriority    `json:"priority"`
	Tags              *[]string                  `jcon:"tags"`
}
type UpsertImsiConfigurationData struct {
	UpdatedBy         string
	ProjectId         uuid.UUID
	StationLocationId uuid.UUID                  `json:"stationLocationId"`
	Imsi              string                     `json:"imsi"`
	PermittedLabel    model.DevicePermittedLabel `json:"permittedLabel"`
	BlacklistPriority model.BlacklistPriority    `json:"priority"`
	Tags              *[]string                  `jcon:"tags"`
}

type CreateImsiConfigurationInput struct {
	StationLocationId graphql.ID                `json:"stationLocationId"`
	Imsi              string                    `json:"imsi"`
	PermittedLabel    enum.DevicePermittedLabel `json:"permittedLabel"`
	BlacklistPriority enum.BlacklistPriority    `json:"priority"`
	Tags              *[]string                 `jcon:"tags"`
}

type CreateImsiConfigurationData struct {
	ProjectId         uuid.UUID `json:"projectId"`
	Imsi              string    `json:"imsi"`
	CreatedBy         string    `json:"createdBy"`
	PermittedLabel    string    `json:"label"`
	BlacklistPriority string    `json:"priority"`
	StationLocationId uuid.UUID `json:"stationLocationId"`
	Tags              *[]string `json:"tags"`
}

type UpdateImsiConfigurationData struct {
	ID                uuid.UUID `json:"id"`
	Imsi              *string   `json:"imsi"`
	ProjectId         uuid.UUID `json:"projectId"`
	UpdatedBy         string    `json:"updatedBy"`
	PermittedLabel    *string   `json:"label"`
	BlacklistPriority *string   `json:"priority"`
	Tags              *[]string `json:"tags"`
}

type UpdateImsiConfigurationInput struct {
	ID                graphql.ID                 `json:"id"`
	Imsi              string                     `json:"imsi"`
	ProjectId         graphql.ID                 `json:"projectId"`
	PermittedLabel    *enum.DevicePermittedLabel `json:"label"`
	BlacklistPriority *enum.BlacklistPriority    `json:"priority"`
	Tags              *[]string                  `json:"tags"`
}

type GetImsiConfigurationsData struct {
	StationLocationId uuid.UUID `json:"stationLocationId"`
	Search            *string   `json:"search"`
	ProjectId         uuid.UUID `json:"projectId"`
	PermittedLabel    *string   `json:"permittedLabel"`
	BlacklistPriority *string   `json:"blacklistPriority"`
	Tags              *[]string `json:"tags"`
	Mnc               *string   `json:"mnc"`
	Mcc               *string   `json:"mcc"`
	Pagination        utils.OffsetPagination
}

type GetImsiConfigurationsInput struct {
	StationLocationId graphql.ID                 `json:"stationLocationId"`
	Search            *string                    `json:"search"`
	PermittedLabel    *enum.DevicePermittedLabel `json:"permittedLabel"`
	BlacklistPriority *enum.BlacklistPriority    `json:"blacklistPriority"`
	Tags              *[]string                  `json:"tags"`
	Mnc               *string                    `json:"mnc"`
	Mcc               *string                    `json:"mcc"`
	Limit             float64                    `json:"limit"`
	Skip              float64                    `json:"skip"`
}

type GetImsiConfigurationByIdInput struct {
	ID graphql.ID `json:"id"`
}

type GetImsiConfigurationByIdData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
}

type DeleteImsiConfigurationInput struct {
	ID graphql.ID `json:"id"`
}

type DeleteImsiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	DeletedBy string    `json:"updatedBy"`
}
