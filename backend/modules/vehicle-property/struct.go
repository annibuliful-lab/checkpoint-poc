package vehicleproperty

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/gql/enum"

	"github.com/google/uuid"
)

type VehicleProperty struct {
	Property string
	Type     model.PropertyType
}

type UpsertVehiclePropertyData struct {
	ProjectId uuid.UUID
	Property  string
	Type      model.PropertyType
}

type GetVehiclePropertiesInput struct {
	Type   enum.PropertyType
	Search *string
	Limit  float64
	Skip   float64
}

type GetVehiclePropertiesData struct {
	ProjectId uuid.UUID
	Type      model.PropertyType
	Search    *string
	Limit     int64
	Skip      int64
}
