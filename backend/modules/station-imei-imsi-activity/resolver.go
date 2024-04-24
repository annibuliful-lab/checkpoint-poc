package stationimeiimsiactivity

import (
	"checkpoint/gql/enum"
	imeiconfiguration "checkpoint/modules/imei-configuration"
	imsiconfiguration "checkpoint/modules/imsi-configuration"
	"context"
	"strings"
	"time"

	"github.com/graph-gophers/graphql-go"
)

type StationImeiImsiActivityResolver struct{}

func (StationImeiImsiActivityResolver) GetStationImeiImsiActivitySummary(ctx context.Context, args StationImeiImsiActivitySummaryData) (*StationImeiImsiActivitySummary, error) {

	return &StationImeiImsiActivitySummary{
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

func (parent StationImeiImsiActivitySummary) Series(ctx context.Context, args StationImeiImsiActivitySummarySerieFilter) (*[]StationImeiImsiActivitySummarySerie, error) {
	var results []StationImeiImsiActivitySummarySerie
	imeiData := StationImeiImsiActivitySummarySerie{
		Label: "IMEI",
		Data:  []int32{10, 41, 35, 51, 49, 62, 69, 91, 148},
	}
	imsiData := StationImeiImsiActivitySummarySerie{
		Label: "IMSI",
		Data:  []int32{10, 8, 35, 7, 49, 6, 69, 8, 148},
	}
	switch {
	case args.Filter != nil && strings.EqualFold(*args.Filter, "IMEI"):
		results = append(results, imeiData)
	case args.Filter != nil && strings.EqualFold(*args.Filter, "IMSI"):
		results = append(results, imsiData)
	default:
		results = append(results, imsiData, imeiData)
	}
	return &results, nil
}

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

func (StationImeiImsiActivityResolver) GetStationImeiImsiActivityById(ctx context.Context, input struct{ ID graphql.ID }) (*StationImeiImsiActivity, error) {
	currentTime := time.Now()
	isoTimeString := currentTime.Format(time.RFC3339)
	return &StationImeiImsiActivity{
		ID:                "Mock-ID",
		ProjectId:         "Mock-ProjectId",
		StationLocationId: "Mock-StationLocationId",
		ArrivalTime:       isoTimeString,
		LicensePlate:      "Mock-LicensePlate",
		PhoneModel:        "Mock-PhoneModel",
		StationSiteName:   "Mock-StationSiteName",
	}, nil
}

func (parent StationImeiImsiActivity) Tags(ctx context.Context) (*[]StationImeiImsiActivityTag, error) {
	tags := []StationImeiImsiActivityTag{{Tag: "Mock-Tag License", Type: enum.GetActivityTagType("LICENSE_PLATE")}}

	return &tags, nil
}

func (parent StationImeiImsiActivity) Imsi(ctx context.Context) (*imsiconfiguration.ImsiConfiguration, error) {
	graphqlTime := graphql.Time{Time: time.Now()}
	imsi := imsiconfiguration.ImsiConfiguration{
		ID:                "Mock-ID",
		ProjectId:         "Mock-ProjectId",
		Imsi:              "Mock-Imsi",
		CreatedBy:         "Mock-CreatedBy",
		CreatedAt:         graphqlTime,
		BlacklistPriority: enum.GetBlacklistPriority("NORMAL"),
		StationLocationId: "Mock-StationLocationId",
		PermittedLabel:    enum.GetDevicePermittedLabel("NONE"),
	}
	return &imsi, nil
}

func (parent StationImeiImsiActivity) Imei(ctx context.Context) (*imeiconfiguration.ImeiConfiguration, error) {
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
