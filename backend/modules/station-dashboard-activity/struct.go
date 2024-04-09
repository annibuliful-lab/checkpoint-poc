package stationdashboardactivity

import (
	"checkpoint/gql/enum"

	"github.com/graph-gophers/graphql-go"
)

type StationDashboardActivityData struct {
	ProjectId graphql.ID `json:"projectId"`
	StationId graphql.ID `json:"stationId"`
	Limit     int32      `json:"limit"`
	Skip      int32      `json:"skip"`
}

type StationDashboardActivity struct {
	ID                graphql.ID `json:"id"`
	ProjectId         graphql.ID `json:"projectId"`
	StationLocationId graphql.ID `json:"stationLocationId"`
	ArrivalTime       string     `json:"arrivalTime"`
	PhoneModel        string     `json:"phoneModel"`
	LicensePlate      string     `json:"licensePlate"`
	StationSiteName   string     `json:"stationSiteName"`
}

type StationDashboardActivityTag struct {
	Type enum.ActivityTagType `json:"type"`
	Tag  string               `json:"tag"`
}

type StationDashboardActivityPicture struct {
	Driver       string
	LicensePlate string
	Front        string
	Back         string
	Side         string
}

type StationDashboardActivityVehicleInfo struct {
	LicensePlate     string
	LicensePlateType string
	VehicleType      string
	StationSiteName  string
	Status           enum.DevicePermittedLabel
	Band             string
	ColorName        string
	ColorCode        string
}
