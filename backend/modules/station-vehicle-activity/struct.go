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

type StationVehicleActivitySummaryData struct {
	StationId  graphql.ID
	GroupBy    string
	CustomDate *graphql.NullTime
}

type StationVehicleActivityTag struct {
	Type enum.StationVehicleActivityTagStatus
	Tag  string
}
type StationVehicleActivityImei struct {
	Total int32
	List  *[]string
}
type StationVehicleActivityImsi struct {
	Total int32
	List  *[]string
}

type StationVehicleActivityColor struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type StationVehicleActivityVehicle struct {
	Type string `json:"type"`
}

type StationVehicleActivity struct {
	ID          string `json:"id"`
	ArrivalTime string `json:"arrivalTime"`
	Brand       string `json:"brand"`
	StationSite string `json:"stationSite"`
	Remark      string `json:"remark"`
}

type StationVehicleActivitySummarySerie struct {
	Label string
	Data  []int32
}

type StationVehicleActivitySummary struct {
	Categories []string
}

type StationVehicleActivityLicensePlate struct {
	Image   *string                   `json:"image"`
	License string                    `json:"license"`
	Type    string                    `json:"type"`
	Status  enum.DevicePermittedLabel `json:"status"`
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
