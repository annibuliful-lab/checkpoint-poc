package stationvehicleactivity

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	"checkpoint/gql/enum"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

var stationVehicleActivityService = StationVehicleActivityService{}

type StationVehicleActivityResolver struct{}

func (StationVehicleActivityResolver) GetStationVehicleActivities(ctx context.Context, args StationVehicleActivityData) (*[]StationVehicleActivity, error) {
	remark := `Mock-Remark`
	return &[]StationVehicleActivity{
		{
			ID:                "Mock-ID",
			ProjectId:         "Mock-ProjectId",
			StationLocationId: "Mock-StationLocationId",
			ArrivalTime:       "Mock-ArrivalTime",
			LicensePlate:      "Mock-licensePlate",
			LicensePlateType:  "Mock-licensePlateType",
			Brand:             "Mock-brand",
			VehicleType:       "Mock-vehicleType",
			Color:             "Mock-color",
			ColorName:         "Mock-colorName",
			Remark:            &remark,
		},
		{
			ID:                "Mock-ID-02",
			ProjectId:         "Mock-ProjectId",
			StationLocationId: "Mock-StationLocationId",
			ArrivalTime:       "Mock-ArrivalTime",
			LicensePlate:      "Mock-licensePlate",
			LicensePlateType:  "Mock-licensePlateType",
			Brand:             "Mock-brand",
			VehicleType:       "Mock-vehicleType",
			Color:             "Mock-color",
			ColorName:         "Mock-colorName",
			Remark:            &remark,
		},
	}, nil
}

func (StationVehicleActivityResolver) GetStationVehicleActivityById(ctx context.Context, input struct{ ID graphql.ID }) (*StationVehicleActivity, error) {
	remark := `Mock-Remark`
	return &StationVehicleActivity{
		ID:                "Mock-ID",
		ProjectId:         "Mock-ProjectId",
		StationLocationId: "Mock-StationLocationId",
		ArrivalTime:       "Mock-ArrivalTime",
		LicensePlate:      "Mock-licensePlate",
		LicensePlateType:  "Mock-licensePlateType",
		Brand:             "Mock-brand",
		VehicleType:       "Mock-vehicleType",
		Color:             "Mock-color",
		ColorName:         "Mock-colorName",
		Remark:            &remark,
	}, nil
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

func (parent StationVehicleActivity) Tags(ctx context.Context) (*[]StationVehicleActivityTag, error) {
	tags := []StationVehicleActivityTag{
		{
			Status: enum.GetStationVehicleActivityTagStatus("IMSI"),
			Tag:    "Mock-Tag Imsi",
		},
		{
			Status: enum.GetStationVehicleActivityTagStatus("IMEI"),
			Tag:    "Mock-Tag Imei",
		},
		{
			Status: enum.GetStationVehicleActivityTagStatus("LICENSE_PLATE"),
			Tag:    "Mock-Tag License plate",
		},
	}

	return &tags, nil
}

func (parent StationVehicleActivity) Imeis(ctx context.Context) (*[]StationVehicleActivityImei, error) {
	imeis := []StationVehicleActivityImei{
		{
			Status: enum.GetRemarkState("WHITELIST"),
			Imei:   "Mock-Imei Danger",
		},
		{
			Status: enum.GetRemarkState("BLACKLIST"),
			Imei:   "Mock-Imei Warning",
		},
		{
			Status: enum.GetRemarkState("PASSED"),
			Imei:   "Mock-Imei",
		},
	}

	return &imeis, nil
}

func (parent StationVehicleActivity) Imsis(ctx context.Context) (*[]StationVehicleActivityImsi, error) {
	imsis := []StationVehicleActivityImsi{
		{
			Status: enum.GetRemarkState("WHITELIST"),
			Imsi:   "Mock-Imsi Danger",
		},
		{
			Status: enum.GetRemarkState("BLACKLIST"),
			Imsi:   "Mock-Imsi Warning",
		},
		{
			Status: enum.GetRemarkState("PASSED"),
			Imsi:   "Mock-Imsi",
		},
	}

	return &imsis, nil
}
