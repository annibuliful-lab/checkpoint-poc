package vehicletargetconfigurationimage

import (
	"checkpoint/gql/enum"

	"github.com/graph-gophers/graphql-go"
)

type VehicleTargetConfigurationImage struct {
	Id                           graphql.ID
	VehicleTargetConfigurationId graphql.ID
	Type                         enum.ImageType
	S3Key                        string
	Url                          string
}
