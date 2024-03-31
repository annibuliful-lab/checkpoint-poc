package mobiledeviceconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/gql/enum"
	"checkpoint/utils"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type GetMobileDeviceConfigurationsInput struct {
	StationLocationId graphql.ID                 `json:"stationLocationId"`
	Search            *string                    `json:"search"`
	PermittedLabel    *enum.DevicePermittedLabel `json:"permittedLabel"`
	BlacklistPriority *enum.BlacklistPriority    `json:"blacklistPriority"`
	Tags              *[]string                  `json:"tags"`
	Limit             float64                    `json:"limit"`
	Skip              float64                    `json:"skip"`
}

type GetMobileDeviceConfigurationsData struct {
	StationLocationId uuid.UUID `json:"stationLocationId"`
	Search            *string   `json:"search"`
	PermittedLabel    *string   `json:"permittedLabel"`
	BlacklistPriority *string   `json:"blacklistPriority"`
	Tags              *[]string `json:"tags"`
	ProjectId         uuid.UUID `json:"projectId"`
	pagination        utils.OffsetPagination
}

type GetMobileDeviceConfigurationInput struct {
	ID graphql.ID `json:"id"`
}

type GetMobileDeviceConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
}

type UpdateMobileDeviceConfigurationInput struct {
	ID                graphql.ID                 `json:"id"`
	Title             *string                    `json:"title"`
	Msisdn            *string                    `json:"msisdn"`
	Imsi              *string                    `json:"referenceImsiConfigurationId"`
	Imei              *string                    `json:"referenceImeiConfigurationId"`
	PermittedLabel    *enum.DevicePermittedLabel `json:"permittedLabel"`
	BlacklistPriority *enum.BlacklistPriority    `json:"blacklistPriority"`
	Tags              *[]string                  `json:"tags"`
}

type UpdateMobileDeviceConfigurationData struct {
	ID                uuid.UUID `json:"id"`
	Title             *string   `json:"title"`
	Msisdn            *string   `json:"msisdn"`
	Imsi              *string   `json:"referenceImsiConfigurationId"`
	Imei              *string   `json:"referenceImeiConfigurationId"`
	PermittedLabel    *string   `json:"permittedLabel"`
	BlacklistPriority *string   `json:"blacklistPriority"`
	ProjectId         uuid.UUID `json:"projectId"`
	StationId         uuid.UUID
	UpdatedBy         string    `json:"updatedBy"`
	Tags              *[]string `json:"tags"`
}

type DeleteMobileDeviceConfigurationInpt struct {
	ID graphql.ID `json:"id"`
}

type DeleteMobileDeviceConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	DeletedBy string    `json:"deletedBy"`
}

type CreateMobileDeviceConfigurationInput struct {
	StationLocationId graphql.ID                `json:"stationLocationId"`
	Title             string                    `json:"title"`
	Msisdn            *string                   `json:"msisdn"`
	Imsi              string                    `json:"referenceImsiConfigurationId"`
	Imei              string                    `json:"referenceImeiConfigurationId"`
	PermittedLabel    enum.DevicePermittedLabel `json:"permittedLabel"`
	BlacklistPriority enum.BlacklistPriority    `json:"blacklistPriority"`
	Tags              *[]string                 `json:"tags"`
}

type CreateMobileDeviceConfigurationData struct {
	StationId         uuid.UUID
	ProjectId         uuid.UUID                  `json:"projectId"`
	StationLocationId uuid.UUID                  `json:"stationLocationId"`
	Title             string                     `json:"title"`
	Msisdn            *string                    `json:"msisdn"`
	Tags              *[]string                  `json:"tags"`
	Imsi              string                     `json:"referenceImsiConfigurationId"`
	Imei              string                     `json:"referenceImeiConfigurationId"`
	PermittedLabel    model.DevicePermittedLabel `json:"permittedLabel"`
	BlacklistPriority model.BlacklistPriority    `json:"blacklistPriority"`
	CreatedBy         string                     `json:"createdBy"`
}

type MobileDeviceConfiguration struct {
	ID                           graphql.ID                `json:"id"`
	StationLocationId            graphql.ID                `json:"stationLocationId"`
	ProjectId                    graphql.ID                `json:"projectId"`
	Title                        string                    `json:"title"`
	ReferenceImsiConfigurationId graphql.ID                `json:"referenceImsiConfigurationId"`
	ReferenceImeiConfigurationId graphql.ID                `json:"referenceImeiConfigurationId"`
	Msisdn                       *string                   `json:"msisdn"`
	PermittedLabel               enum.DevicePermittedLabel `json:"permittedLabel"`
	BlacklistPriority            enum.BlacklistPriority    `json:"blacklistPriority"`
	CreatedBy                    graphql.ID                `json:"createdBy"`
	UpdatedBy                    *graphql.NullID           `json:"updatedBy"`
	CreatedAt                    graphql.Time              `json:"createdAt"`
	UpdatedAt                    *graphql.NullTime         `json:"updatedAt"`
}
