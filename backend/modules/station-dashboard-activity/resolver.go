package stationdashboardactivity

import (
	"checkpoint/gql/enum"
	imeiconfiguration "checkpoint/modules/imei-configuration"
	imsiconfiguration "checkpoint/modules/imsi-configuration"
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
			PhoneModel:        "Mock-PhoneModel",
			StationSiteName:   "Mock-StationSiteName",
		},
		{
			ID:                "Mock-ID-02",
			ProjectId:         "Mock-ProjectId",
			StationLocationId: "Mock-StationLocationId",
			ArrivalTime:       isoTimeString,
			LicensePlate:      "Mock-LicensePlate",
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
		PhoneModel:        "Mock-PhoneModel",
		StationSiteName:   "Mock-StationSiteName",
	}, nil
}

func (parent StationDashboardActivity) Imsi(ctx context.Context) (*imsiconfiguration.ImsiConfiguration, error) {
	graphqlTime := graphql.Time{Time: time.Now()}
	imsi := imsiconfiguration.ImsiConfiguration{
		ID:                "Mock-ID",
		ProjectId:         "Mock-ProjectId",
		Imsi:              "Mock-Imsi",
		CreatedBy:         "Mock-CreatedBy",
		CreatedAt:         graphqlTime,
		BlacklistPriority: enum.GetBlacklistPriority("NORMAL"),
		StationLocationId: "Mock-StationLocationId",
		PermittedLabel:    enum.GetDevicePermittedLabel("BLACKLIST"),
	}
	return &imsi, nil
}

func (parent StationDashboardActivity) Imei(ctx context.Context) (*imeiconfiguration.ImeiConfiguration, error) {
	graphqlTime := graphql.Time{Time: time.Now()}
	imei := imeiconfiguration.ImeiConfiguration{
		ID:                "Mock-ID",
		ProjectId:         "Mock-ProjectId",
		Imei:              "Mock-Imei",
		CreatedBy:         "Mock-CreatedBy",
		CreatedAt:         graphqlTime,
		BlacklistPriority: enum.GetBlacklistPriority("NORMAL"),
		StationLocationId: "Mock-StationLocationId",
		PermittedLabel:    enum.GetDevicePermittedLabel("NONE"),
	}
	return &imei, nil
}

func (parent StationDashboardActivity) Tags(ctx context.Context) (*[]StationDashboardActivityTag, error) {
	tags := []StationDashboardActivityTag{{Tag: "tag", Type: enum.GetActivityTagType("LICENSE_PLATE")}}

	return &tags, nil
}

func (parent StationDashboardActivity) Pictures(ctx context.Context) (*StationDashboardActivityPicture, error) {
	pictures := StationDashboardActivityPicture{
		Driver:       "Mock-Driver",
		LicensePlate: "Mock-LicensePlate",
		Front:        "Mock-Front",
		Back:         "Mock-Back",
		Side:         "Mock-Side",
	}

	return &pictures, nil
}

func (parent StationDashboardActivity) VehicleInfo(ctx context.Context) (*StationDashboardActivityVehicleInfo, error) {
	vehicleInfo := StationDashboardActivityVehicleInfo{
		Status:           enum.GetDevicePermittedLabel("NONE"),
		LicensePlate:     "Mock-LicensePlate",
		LicensePlateType: "Mock-LicensePlateType",
		VehicleType:      "Mock-VehicleType",
		StationSiteName:  "Mock-StationSiteName",
		Band:             "Mock-Band",
		ColorName:        "Mock-ColorName",
		ColorCode:        "Mock-ColorCode",
	}

	return &vehicleInfo, nil
}
