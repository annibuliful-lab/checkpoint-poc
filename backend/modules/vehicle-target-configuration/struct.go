package vehicletargetconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/gql/enum"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type VehicleTargetConfiguration struct {
	StationLocationId graphql.ID
	ID                graphql.ID
	ProjectId         graphql.ID
	Color             string
	Brand             string
	Type              string
	Prefix            string
	Number            string
	Province          string
	Country           *string
	PermittedLabel    enum.DevicePermittedLabel
	BlacklistPriority enum.BlacklistPriority
}

type UpsertImageS3KeyInput struct {
	Id    *graphql.ID
	S3Key string
	Type  enum.ImageType
}

type CreateVehicleTargetConfigurationInput struct {
	StationLocationId graphql.ID
	Color             string
	Brand             string
	Type              string
	Prefix            string
	Number            string
	Province          string
	Country           *string
	Images            *[]string
	PermittedLabel    enum.DevicePermittedLabel
	BlacklistPriority enum.BlacklistPriority
	Tags              *[]string
	ImageS3Keys       *[]UpsertImageS3KeyInput
}

type CreateVehicleTargetConfigurationData struct {
	StationLocationId uuid.UUID
	ProjectId         uuid.UUID
	Color             string
	Brand             string
	Type              string
	CreatedBy         string
	Prefix            string
	Number            string
	Province          string
	Country           *string
	Images            *[]string
	PermittedLabel    model.DevicePermittedLabel
	BlacklistPriority model.BlacklistPriority
	Tags              *[]string
}

type UpdateVehicleTargetConfigurationInput struct {
	ID                graphql.ID
	Color             *string
	Brand             *string
	Prefix            *string
	Number            *string
	Province          *string
	Type              *string
	Country           *string
	Images            *[]string
	PermittedLabel    *enum.DevicePermittedLabel
	BlacklistPriority *enum.BlacklistPriority
	Tags              *[]string
	ImageS3Keys       *[]UpsertImageS3KeyInput
}

type UpdateVehicleTargetConfigurationData struct {
	ID                uuid.UUID
	ProjectId         uuid.UUID
	UpdatedBy         string
	Color             *string
	Brand             *string
	Model             *string
	Prefix            *string
	Number            *string
	Province          *string
	Type              *string
	Country           *string
	Images            *[]string
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
	StationLocationId *graphql.ID
	Search            *string
	Type              *string
	PermittedLabel    *enum.DevicePermittedLabel
	BlacklistPriority *enum.BlacklistPriority
	Tags              *[]string
	Limit             float64
	Skip              float64
}

type GetVehicleTargetsConfigurationData struct {
	StationLocationId *uuid.UUID
	ProjectId         uuid.UUID
	Type              *string
	Search            *string
	PermittedLabel    *model.DevicePermittedLabel
	BlacklistPriority *model.BlacklistPriority
	Tags              *[]string
	Limit             int64
	Skip              int64
}
