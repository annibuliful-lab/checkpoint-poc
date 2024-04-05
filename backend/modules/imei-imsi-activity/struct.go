package imeiimsiactivity

import (
	"checkpoint/gql/enum"

	"github.com/graph-gophers/graphql-go"
)

type ImeiImsiActivityData struct {
	ProjectId graphql.ID `json:"projectId"`
	StationId graphql.ID `json:"stationId"`
	Limit     int32      `json:"limit"`
	Skip      int32      `json:"skip"`
}

type ImeiImsiActivity struct {
	ID                graphql.ID `json:"id"`
	ProjectId         graphql.ID `json:"projectId"`
	StationLocationId graphql.ID `json:"stationLocationId"`
	ArrivalTime       string     `json:"arrivalTime"`
	Imei              string     `json:"imei"`
	Imsi              string     `json:"imsi"`
	PhoneModel        string     `json:"phoneModel"`
	LicensePlate      string     `json:"licensePlate"`
	StationSiteName   string     `json:"stationSiteName"`
}

type ImeiImsiActivityTag struct {
	Type enum.ActivityTagType `json:"type"`
	Tag  string               `json:"tag"`
}
