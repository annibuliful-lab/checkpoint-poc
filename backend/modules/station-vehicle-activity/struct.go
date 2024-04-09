package stationvehicleactivity

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/gql/enum"
	"checkpoint/utils/graphql_utils"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type StationVehicleActivityData struct {
	ProjectId graphql.ID `json:"projectId"`
	StationId graphql.ID `json:"stationId"`
	Limit     int32      `json:"limit"`
	Skip      int32      `json:"skip"`
}

type StationVehicleActivityTag struct {
	Status enum.StationVehicleActivityTagStatus
	Tag    string
}
type StationVehicleActivityImei struct {
	Status enum.RemarkState
	Imei   string
}
type StationVehicleActivityImsi struct {
	Status enum.RemarkState
	Imsi   string
}
type StationVehicleActivity struct {
	ID                graphql.ID `json:"id"`
	ProjectId         graphql.ID `json:"projectId"`
	StationLocationId graphql.ID `json:"stationLocationId"`
	ArrivalTime       string     `json:"arrivalTime"`
	LicensePlate      string     `json:"licensePlate"`
	LicensePlateType  string     `json:"licensePlateType"`
	Brand             string     `json:"brand"`
	VehicleType       string     `json:"vehicleType"`
	Color             string     `json:"color"`
	ColorName         string     `json:"colorName"`
	Remark            *string    `json:"remark"`
}

type StationVehicleActivityConnection struct {
	PageInfo graphql_utils.PageInfo   `json:"pageInfo"`
	Edges    []StationVehicleActivity `json:"edges"`
}

type UpdateStationVehicleActivityInput struct {
	ID     graphql.ID
	Brand  *string
	Color  *string
	Model  *string
	Status *enum.RemarkState
}

type UpdateStationVehicleActivityData struct {
	ID        uuid.UUID
	ProjectId uuid.UUID
	UpdatedBy string
	Brand     *string
	Color     *string
	Model     *string
	Status    *enum.RemarkState
}

type CreateStationVehicleActivityInput struct {
	Brand  string
	Color  string
	Model  string
	Status *enum.RemarkState
}

type CreateStationVehicleActivityData struct {
	ProjectId         uuid.UUID
	StationLocationId uuid.UUID
	Brand             string
	Color             string
	Model             string
	Status            *model.RemarkState
	CreatedBy         string
}
