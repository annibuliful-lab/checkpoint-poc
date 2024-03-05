package stationhealthcheck

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/gql/enum"
	"time"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type StationLocationHealthCheckActivity struct {
	ID            graphql.ID
	StationId     graphql.ID
	StationStatus enum.StationStatus
	StartDatetime graphql.Time
	EndDatetime   *graphql.Time
}

type CreateStationHealthCheckActivityInput struct {
	StationId     graphql.ID
	StationStatus enum.StationStatus
	StartDatetime graphql.Time
	EndDatetime   *graphql.NullTime
}

type CreateStationHealthCheckActivityData struct {
	StationId     uuid.UUID
	StationStatus model.StationStatus
	StartDatetime time.Time
	EndDatetime   *time.Time
	CreatedBy     string
}

type UpdateStationHealthCheckActivityInput struct {
	ID            graphql.ID
	StationStatus *enum.StationStatus
	StartDatetime *graphql.Time
	EndDatetime   *graphql.Time
}

type UpdateStationHealthCheckActivityData struct {
	ID            uuid.UUID
	StationStatus *model.StationStatus
	StartDatetime *time.Time
	EndDatetime   *time.Time
	UpdatedBy     string
}

type GetStationHealthCheckActivitiesInput struct {
	StationId     graphql.ID
	StationStatus *enum.StationStatus
	StartDatetime *graphql.Time
	EndDatetime   *graphql.Time
	Limit         float64
	Skip          float64
}

type GetStationHealthCheckActivitiesData struct {
	StationId     uuid.UUID
	StationStatus *model.StationStatus
	StartDatetime *time.Time
	EndDatetime   *time.Time
	Limit         int64
	Skip          int64
}
