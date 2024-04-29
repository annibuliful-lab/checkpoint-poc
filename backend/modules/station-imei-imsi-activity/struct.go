package stationimeiimsiactivity

import (
	"checkpoint/gql/enum"

	"github.com/graph-gophers/graphql-go"
)

type StationImeiImsiActivityData struct {
	ProjectId graphql.ID `json:"projectId"`
	StationId graphql.ID `json:"stationId"`
	Limit     int32      `json:"limit"`
	Skip      int32      `json:"skip"`
}

type StationImeiImsiActivity struct {
	ID                graphql.ID `json:"id"`
	ProjectId         graphql.ID `json:"projectId"`
	StationLocationId graphql.ID `json:"stationLocationId"`
	ArrivalTime       string     `json:"arrivalTime"`
	PhoneModel        string     `json:"phoneModel"`
	LicensePlate      string     `json:"licensePlate"`
	StationSiteName   string     `json:"stationSiteName"`
}

type StationImeiImsiActivityTag struct {
	Type enum.ActivityTagType `json:"type"`
	Tag  string               `json:"tag"`
}

type StationImeiImsiActivitySummaryData struct {
	StationId  graphql.ID
	GroupBy    string
	CustomDate *graphql.NullTime
}

type StationImeiImsiActivitySummarySerieFilter struct {
	Filter *string
}

type StationImeiImsiActivitySummarySerie struct {
	Label string
	Data  []int32
}

type StationImeiImsiActivitySummary struct {
	Categories []string
}
