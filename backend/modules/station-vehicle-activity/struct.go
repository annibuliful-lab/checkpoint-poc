package stationvehicleactivity

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/gql/enum"
	"checkpoint/utils/graphql_utils"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type StationVehicleActivity struct {
	ID                graphql.ID
	ProjectId         graphql.ID
	StationLocationId graphql.ID
	Brand             string
	Color             string
	Model             string
	Status            enum.RemarkState
}

type StationVehicleActivityConnection struct {
	PageInfo graphql_utils.PageInfo
	Edges    []StationVehicleActivity
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
