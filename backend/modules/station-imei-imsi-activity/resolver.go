package stationimeiimsiactivity

import (
	"checkpoint/gql/enum"
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
)

type StationImeiImsiActivityResolver struct{}

func (StationImeiImsiActivityResolver) GetStationImeiImsiActivities(ctx context.Context, args StationImeiImsiActivityData) (*[]StationImeiImsiActivity, error) {
	currentTime := time.Now()
	isoTimeString := currentTime.Format(time.RFC3339)
	return &[]StationImeiImsiActivity{
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

func (StationImeiImsiActivityResolver) GetStationImeiImsiActivityById(ctx context.Context, input struct{ ID graphql.ID }) (*StationImeiImsiActivity, error) {
	currentTime := time.Now()
	isoTimeString := currentTime.Format(time.RFC3339)
	return &StationImeiImsiActivity{
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

func (parent StationImeiImsiActivity) Tags(ctx context.Context) (*[]StationImeiImsiActivityTag, error) {
	tags := []StationImeiImsiActivityTag{{Tag: "Mock-Tag License", Type: enum.GetActivityTagType("LICENSE_PLATE")}}

	return &tags, nil
}
