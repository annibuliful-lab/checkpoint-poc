package vehicletargetconfigurationimage

import (
	"checkpoint/gql/enum"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type VehicleTargetConfigurationImage struct {
	Id                           graphql.ID
	VehicleTargetConfigurationId graphql.ID
	Type                         enum.ImageType
	S3Key                        string
	Url                          string
}

type CreateVehicleTargetConfigurationImageData struct {
	VehicleTargetConfigurationId uuid.UUID
	Type                         enum.ImageType
	S3Key                        string
	CreatedBy                    string
}

type UpdateVehicleTargetConfigurationImageData struct {
	Id        uuid.UUID
	Type      enum.ImageType
	S3Key     string
	UpdatedBy string
}

type DeleteVehicleTargetConfigurationImageInput struct {
	Id graphql.ID
}

type DeleteVehicleTargetConfigurationImageData struct {
	Id        uuid.UUID
	DeletedBy string
}
