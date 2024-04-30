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

type StationVehicleActivityResolver struct {
	stationVehicleActivityEvent      chan *StationVehicleActivity
	stationVehicleActivitySubscriber chan *StationVehicleActivitySubscriber
}

func (StationVehicleActivityResolver) GetStationVehicleActivitySummary(ctx context.Context, args StationVehicleActivitySummaryData) (*StationVehicleActivitySummary, error) {

	return &StationVehicleActivitySummary{
		Categories: []string{
			"Jan",
			"Feb",
			"Mar",
			"Apr",
			"May",
			"Jun",
			"Jul",
			"Aug",
			"Sep"},
	}, nil
}

func (StationVehicleActivitySummary) Series(ctx context.Context) (*[]StationVehicleActivitySummarySerie, error) {

	return &[]StationVehicleActivitySummarySerie{{
		Label: "Vehicle count",
		Data:  []int32{10, 41, 35, 51, 49, 62, 69, 91, 148},
	}}, nil
}

func (StationVehicleActivityResolver) GetStationVehicleActivities(ctx context.Context, args StationVehicleActivityData) (*[]StationVehicleActivity, error) {

	return &[]StationVehicleActivity{{
		ID:          "Mock-ID-1",
		ArrivalTime: "Thu Apr 25 2024 13:39:46 GMT+0700 (Indochina Time)",
		Brand:       "Honda",
		StationSite: "1",
		Remark:      "Remark-note-1",
	}, {
		ID:          "Mock-ID-2",
		ArrivalTime: "Thu Apr 25 2024 13:39:46 GMT+0700 (Indochina Time)",
		Brand:       "Yamaha",
		StationSite: "1",
		Remark:      "Remark-note-2",
	}}, nil
}

func (StationVehicleActivityResolver) GetStationVehicleActivityById(ctx context.Context, input struct{ ID graphql.ID }) (*StationVehicleActivity, error) {
	return &StationVehicleActivity{}, nil
}

func (StationVehicleActivityResolver) UpdateStationVehicleActivity(ctx context.Context, input UpdateStationVehicleActivityInput) (*StationVehicleActivity, error) {

	return &StationVehicleActivity{}, nil
}

func (r *StationVehicleActivityResolver) CreateStationVehicleActivity(ctx context.Context, input CreateStationVehicleActivityInput) (*StationVehicleActivity, error) {

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

	go func() {
		r.stationVehicleActivityEvent <- vehicleActivity
	}()

	return vehicleActivity, nil
}

func (parent StationVehicleActivity) LicensePlate(ctx context.Context) (*StationVehicleActivityLicensePlate, error) {

	imageDB := map[string]string{
		"Mock-ID-1": "https://global.discourse-cdn.com/freecodecamp/original/3X/e/2/e27c45a168d1c6da0630abcb30d1e7f0d49a4bc6.png",
		"Mock-ID-2": "https://example.com/image2.jpg",
	}

	findImageByStationVehicleActivityId := imageDB[parent.ID]

	licensePlate := StationVehicleActivityLicensePlate{
		Image:   &findImageByStationVehicleActivityId,
		License: "1กท 7777",
		Type:    "12",
		Status:  enum.DevicePermittedLabel(2),
	}
	return &licensePlate, nil
}

func (parent StationVehicleActivity) Vehicle(ctx context.Context) (*StationVehicleActivityVehicle, error) {
	vehicle := StationVehicleActivityVehicle{
		Type: "จักรยานยนต์",
	}

	return &vehicle, nil
}

func (parent StationVehicleActivity) Color(ctx context.Context) (*StationVehicleActivityColor, error) {
	color := StationVehicleActivityColor{
		Name: "aaaa",
		Code: "bbbb",
	}

	return &color, nil
}

func (parent StationVehicleActivity) Tags(ctx context.Context) (*[]StationVehicleActivityTag, error) {
	tags := []StationVehicleActivityTag{{Tag: "Tag-1", Type: enum.GetStationVehicleActivityTagStatus("LICENSE_PLATE"), StationVehicleActivityId: "Mock-ID-1"}, {Tag: "Tag-2", Type: enum.GetStationVehicleActivityTagStatus("LICENSE_PLATE"), StationVehicleActivityId: "Mock-ID-2"}, {Tag: "Tag-3", Type: enum.GetStationVehicleActivityTagStatus("LICENSE_PLATE"), StationVehicleActivityId: "Mock-ID-2"}}

	result := []StationVehicleActivityTag{}

	for _, tag := range tags {
		if tag.StationVehicleActivityId == parent.ID {
			result = append(result, tag)
		}
	}

	return &result, nil
}

func (parent StationVehicleActivity) Imei(ctx context.Context) (*StationVehicleActivityImei, error) {
	imei := StationVehicleActivityImei{}

	return &imei, nil
}

func (parent StationVehicleActivity) Imsi(ctx context.Context) (*StationVehicleActivityImsi, error) {
	imsi := StationVehicleActivityImsi{}

	return &imsi, nil
}
