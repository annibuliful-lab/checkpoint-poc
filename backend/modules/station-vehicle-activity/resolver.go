package stationvehicleactivity

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

var stationVehicleActivityService = StationVehicleActivityService{}

type StationVehicleActivityResolver struct{}

func (StationVehicleActivityResolver) GetStationVehicleActivities(ctx context.Context, args StationVehicleActivityData) (*[]StationVehicleActivity, error) {

	return &[]StationVehicleActivity{}, nil
}

func (StationVehicleActivityResolver) GetStationVehicleActivityById(ctx context.Context, input struct{ ID graphql.ID }) (*StationVehicleActivity, error) {
	return &StationVehicleActivity{}, nil
}

func (StationVehicleActivityResolver) UpdateStationVehicleActivity(ctx context.Context, input UpdateStationVehicleActivityInput) (*StationVehicleActivity, error) {

	return &StationVehicleActivity{}, nil
}

func (StationVehicleActivityResolver) CreateStationVehicleActivity(ctx context.Context, input CreateStationVehicleActivityInput) (*StationVehicleActivity, error) {

	stationApiAccess := auth.GetStationAuthorizationContext(ctx)

	fields := CreateStationVehicleActivityData{
		ProjectId:         uuid.MustParse(stationApiAccess.ProjectId),
		StationLocationId: uuid.MustParse(stationApiAccess.StationId),
		CreatedBy:         stationApiAccess.DeviceId,
		Brand:             input.Brand,
		Color:             input.Color,
		Model:             input.Model,
	}

	if input.Status != nil {
		status := model.RemarkState(input.Status.String())
		fields.Status = &status
	}

	vehicleActivity, err := stationVehicleActivityService.Create(fields)

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return vehicleActivity, nil
}

func (parent StationVehicleActivity) LicensePlate(ctx context.Context) (*StationVehicleActivityLicensePlate, error) {
	licensePlate := StationVehicleActivityLicensePlate{}

	return &licensePlate, nil
}

func (parent StationVehicleActivity) Vehicle(ctx context.Context) (*StationVehicleActivityVehicle, error) {
	vehicle := StationVehicleActivityVehicle{}

	return &vehicle, nil
}

func (parent StationVehicleActivity) Color(ctx context.Context) (*StationVehicleActivityColor, error) {
	color := StationVehicleActivityColor{}

	return &color, nil
}

func (parent StationVehicleActivity) Tags(ctx context.Context) (*[]StationVehicleActivityTag, error) {
	color := []StationVehicleActivityTag{}

	return &color, nil
}

func (parent StationVehicleActivity) Imei(ctx context.Context) (*StationVehicleActivityImei, error) {
	imei := StationVehicleActivityImei{}

	return &imei, nil
}

func (parent StationVehicleActivity) Imsi(ctx context.Context) (*StationVehicleActivityImsi, error) {
	imsi := StationVehicleActivityImsi{}

	return &imsi, nil
}
