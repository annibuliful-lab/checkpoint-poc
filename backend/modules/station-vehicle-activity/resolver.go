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

func (StationVehicleActivityResolver) GetStationVehicleActivityById(ctx context.Context, input struct{ ID graphql.ID }) (*StationVehicleActivity, error) {
	return &StationVehicleActivity{}, nil
}

func (StationVehicleActivityResolver) UpdateStationVehicleActivity(ctx context.Context, input UpdateStationVehicleActivityInput) (*StationVehicleActivity, error) {
	// authorization := auth.GetAuthorizationContext(ctx)

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
