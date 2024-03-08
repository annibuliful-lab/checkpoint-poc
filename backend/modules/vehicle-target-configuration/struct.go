package vehicletargetconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/gql/enum"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type VehicleTargetConfiguration struct {
	ID                graphql.ID
	ProjectId         graphql.ID
	Prefix            string
	Number            string
	Province          string
	Type              string
	Country           *string
	PermittedLabel    enum.DevicePermittedLabel
	BlacklistPriority enum.BlacklistPriority
}

type CreateVehicleTargetConfigurationInput struct {
	Prefix            string
	Number            string
	Province          string
	Type              string
	Country           *string
	PermittedLabel    enum.DevicePermittedLabel
	BlacklistPriority enum.BlacklistPriority
	Tags              *[]string
}

type CreateVehicleTargetConfigurationData struct {
	ProjectId         uuid.UUID
	CreatedBy         string
	Prefix            string
	Number            string
	Province          string
	Type              string
	Country           *string
	PermittedLabel    model.DevicePermittedLabel
	BlacklistPriority model.BlacklistPriority
	Tags              *[]string
}

type UpdateVehicleTargetConfigurationInput struct {
	ID                graphql.ID
	Prefix            *string
	Number            *string
	Province          *string
	Type              *string
	Country           *string
	PermittedLabel    *enum.DevicePermittedLabel
	BlacklistPriority *enum.BlacklistPriority
	Tags              *[]string
}

type UpdateVehicleTargetConfigurationData struct {
	ID                uuid.UUID
	ProjectId         uuid.UUID
	UpdatedBy         string
	Prefix            *string
	Number            *string
	Province          *string
	Type              *string
	Country           *string
	PermittedLabel    *model.DevicePermittedLabel
	BlacklistPriority *model.BlacklistPriority
	Tags              *[]string
}

type DeleteVehicleTargetConfigurationInput struct {
	ID graphql.ID
}

type DeleteVehicleTargetConfigurationData struct {
	ID        uuid.UUID
	ProjectId uuid.UUID
	DeletedBy string
}

type GetVehicleTargetConfigurationInput struct {
	ID graphql.ID
}

type GetVehicleTargetConfigurationData struct {
	ID        uuid.UUID
	ProjectId uuid.UUID
}

type GetVehicleTargetsConfigurationInput struct {
	Search            *string
	Type              *string
	PermittedLabel    *enum.DevicePermittedLabel
	BlacklistPriority *enum.BlacklistPriority
	Tags              *[]string
	Limit             float64
	Skip              float64
}

type GetVehicleTargetsConfigurationData struct {
	ProjectId         uuid.UUID
	Type              *string
	Search            *string
	PermittedLabel    *model.DevicePermittedLabel
	BlacklistPriority *model.BlacklistPriority
	Tags              *[]string
	Limit             int64
	Skip              int64
}
