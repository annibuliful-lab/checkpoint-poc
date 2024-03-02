package stationdevice

import (
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type StationDevice struct {
	ID                graphql.ID
	StationLocationId graphql.ID
	Title             string
	SoftwareVersion   *string
	HardwareVersion   *string
}

type CreateStationDeviceInput struct {
	StationLocationId graphql.ID
	Title             string
	SoftwareVersion   *string
	HardwareVersion   *string
}

type CreateStationDeviceData struct {
	StationLocationId uuid.UUID
	Title             string
	SoftwareVersion   *string
	HardwareVersion   *string
	CreatedBy         string
}

type DeleteStationDeviceInput struct {
	ID graphql.ID
}

type DeleteStationDeviceData struct {
	ID        uuid.UUID
	DeletedBy string
}

type UpdateStationDeviceInput struct {
	ID              graphql.ID
	Title           *string
	SoftwareVersion *string
	HardwareVersion *string
}

type UpdateStationDeviceData struct {
	ID              uuid.UUID
	Title           *string
	SoftwareVersion *string
	HardwareVersion *string
	UpdatedBy       string
}

type GetStationDeviceInput struct {
	ID graphql.ID
}

type GetStationDeviceData struct {
	ID uuid.UUID
}

type GetStationDevicesInput struct {
	StationLocationId graphql.ID
	Search            *string
	Limit             float64
	Skip              float64
}

type GetStationDevicesData struct {
	StationLocationId uuid.UUID
	Search            *string
	Limit             int64
	Skip              int64
}
