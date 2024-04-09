package stationdashboardactivity

import (
	"checkpoint/gql/enum"
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
)

type StationDashboardActivityResolver struct{}

func (StationDashboardActivityResolver) GetStationDashboardActivities(ctx context.Context, args StationDashboardActivityData) (*[]StationDashboardActivity, error) {
	currentTime := time.Now()
	isoTimeString := currentTime.Format(time.RFC3339)
	return &[]StationDashboardActivity{
		{
			ID:                "Mock-ID",
			ProjectId:         "Mock-ProjectId",
			StationLocationId: "Mock-StationLocationId",
			ArrivalTime:       isoTimeString,
			LicensePlate:      "Mock-LicensePlate",
			Imei:              "Mock-Imei",
			Imsi:              "Mock-Imsi",
			PhoneModel:        "Mock-PhoneModel",
			StationSiteName:   "Mock-StationSiteName",
		},
		{
			ID:                "Mock-ID-02",
			ProjectId:         "Mock-ProjectId",
			StationLocationId: "Mock-StationLocationId",
			ArrivalTime:       isoTimeString,
			LicensePlate:      "Mock-LicensePlate",
			Imei:              "Mock-Imei",
			Imsi:              "Mock-Imsi",
			PhoneModel:        "Mock-PhoneModel",
			StationSiteName:   "Mock-StationSiteName",
		},
	}, nil
}

func (StationDashboardActivityResolver) GetStationDashboardActivityById(ctx context.Context, input struct{ ID graphql.ID }) (*StationDashboardActivity, error) {
	currentTime := time.Now()
	isoTimeString := currentTime.Format(time.RFC3339)
	return &StationDashboardActivity{
		ID:                "Mock-ID",
		ProjectId:         "Mock-ProjectId",
		StationLocationId: "Mock-StationLocationId",
		ArrivalTime:       isoTimeString,
		LicensePlate:      "Mock-LicensePlate",
		Imei:              "Mock-Imei",
		Imsi:              "Mock-Imsi",
		PhoneModel:        "Mock-PhoneModel",
		StationSiteName:   "Mock-StationSiteName",
	}, nil
}

func (parent StationDashboardActivity) Tags(ctx context.Context) (*[]StationDashboardActivityTag, error) {
	tags := []StationDashboardActivityTag{{Tag: "Mock-Tag License", Type: enum.GetActivityTagType("LICENSE_PLATE")}}

	return &tags, nil
}

func (parent StationDashboardActivity) Pictures(ctx context.Context) (*StationDashboardActivityPicture, error) {
	pictures := StationDashboardActivityPicture{
		Driver:       "",
		LicensePlate: "",
		Front:        "",
		Back:         "",
		Side:         "",
	}

	return &pictures, nil
}

func (parent StationDashboardActivity) VehicleInfo(ctx context.Context) (*StationDashboardActivityVehicleInfo, error) {
	vehicleInfo := StationDashboardActivityVehicleInfo{
		Status:           enum.GetDevicePermittedLabel("NONE"),
		LicensePlate:     "",
		LicensePlateType: "",
		VehicleType:      "",
		StationSiteName:  "",
		Band:             "",
		ColorName:        "",
		ColorCode:        "",
	}

	return &vehicleInfo, nil
}
