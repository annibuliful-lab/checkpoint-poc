package imeiimsiactivity

import (
	"checkpoint/gql/enum"
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
)

type ImeiImsiActivityResolver struct{}

func (ImeiImsiActivityResolver) GetImeiImsiActivities(ctx context.Context, args ImeiImsiActivityData) (*[]ImeiImsiActivity, error) {
	currentTime := time.Now()
	isoTimeString := currentTime.Format(time.RFC3339)
	return &[]ImeiImsiActivity{
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

func (ImeiImsiActivityResolver) GetImeiImsiActivityById(ctx context.Context, input struct{ ID graphql.ID }) (*ImeiImsiActivity, error) {
	currentTime := time.Now()
	isoTimeString := currentTime.Format(time.RFC3339)
	return &ImeiImsiActivity{
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

func (parent ImeiImsiActivity) Tags(ctx context.Context) (*[]ImeiImsiActivityTag, error) {
	tags := []ImeiImsiActivityTag{{Tag: "Mock-Tag License", Type: enum.GetActivityTagType("LICENSE_PLATE")}}

	return &tags, nil
}
