package vehiclelicenseplate

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/gql/enum"

	"github.com/graph-gophers/graphql-go"
)

type VehicleLicensePlate struct {
	StationVehicleActivityId graphql.ID
	Prefix                   string
	Number                   string
	Province                 string
	Type                     string
	Country                  *string
	Accuracy                 float64
	S3Key                    *string
	ImageUrl                 *string
	PermittedLabel           enum.DevicePermittedLabel
	BlacklistPriority        enum.BlacklistPriority
}

type CreateVehicleLicensePlateInput struct {
	StationVehicleActivityId graphql.ID
	Prefix                   string
	Number                   string
	Province                 string
	Type                     string
	Country                  *string
	Accuracy                 float64
	S3Key                    *string
	ImageUrl                 *string
	PermittedLabel           enum.DevicePermittedLabel
	BlacklistPriority        enum.BlacklistPriority
}

type CreateVehicleLicensePlateData struct {
	StationVehicleActivityId graphql.ID
	Prefix                   string
	Number                   string
	Province                 string
	Type                     string
	Country                  *string
	Accuracy                 float64
	S3Key                    *string
	CreatedBy                string
	PermittedLabel           model.DevicePermittedLabel
	BlacklistPriority        model.BlacklistPriority
}
