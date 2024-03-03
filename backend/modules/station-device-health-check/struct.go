package stationdevicehealthcheck

import (
	"checkpoint/gql/enum"
	"time"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type StationDeviceHealthCheckActivity struct {
	ID              graphql.ID
	StationDeviceId graphql.ID
	Status          enum.DeviceStatus
	Issue           *string
	ActivityTime    graphql.Time
}

type CreateStationDeviceHealthCheckActivityInput struct {
	StationDeviceId graphql.ID
	Status          enum.DeviceStatus
	Issue           *string
	ActivityTime    graphql.Time
}

type CreateStationDeviceHealthCheckActivityData struct {
	StationDeviceId uuid.UUID
	Status          enum.DeviceStatus
	Issue           *string
	ActivityTime    time.Time
}

type GetStationDeviceHealthCheckActivitiesInput struct {
	StationDeviceId graphql.ID
	Status          *enum.DeviceStatus
	StartDate       graphql.NullTime
	EndDate         graphql.NullTime
	Limit           float64
	Skip            float64
}

type GetStationDeviceHealthCheckActivitiesData struct {
	StationDeviceId uuid.UUID
	Status          *enum.DeviceStatus
	StartDate       *time.Time
	EndDate         *time.Time
	Limit           int64
	Skip            int64
}
