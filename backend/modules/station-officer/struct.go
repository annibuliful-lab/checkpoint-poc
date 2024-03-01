package stationofficer

import (
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type StationOfficer struct {
	ID                graphql.ID `json:"id"`
	StationLocationId graphql.ID `json:"stationLocationId"`
	Firstname         string     `json:"firstname"`
	Lastname          string     `json:"lastname"`
	Msisdn            string     `json:"msisdn"`
}

type UpdateStationOfficerInput struct {
	ID        graphql.ID `json:"id"`
	Firstname *string    `json:"firstname"`
	Lastname  *string    `json:"lastname"`
	Msisdn    *string    `json:"msisdn"`
}

type GetStationOfficersInput struct {
	StationLocationId graphql.ID `json:"stationlocationId"`
	Search            *string    `json:"search"`
	Skip              float64    `json:"skip"`
	Limit             float64    `json:"limit"`
}

type GetStationOfficersData struct {
	StationLocationId uuid.UUID `json:"stationlocationId"`
	Search            *string   `json:"search"`
	Skip              int64     `json:"skip"`
	Limit             int64     `json:"limit"`
}

type GetStationOfficerInput struct {
	ID graphql.ID `json:"id"`
}

type GetStationOfficerData struct {
	ID uuid.UUID `json:"id"`
}

type DeleteStationOfficerInput struct {
	ID graphql.ID `json:"id"`
}

type DeleteStationOfficerData struct {
	ID        uuid.UUID `json:"id"`
	DeletedBy string    `json:"deletedBy"`
}

type UpdateStationOfficerData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	Firstname *string   `json:"firstname"`
	Lastname  *string   `json:"lastname"`
	Msisdn    *string   `json:"msisdn"`
	UpdatedBy string    `json:"updatedBy"`
}

type CreateStationOfficerInput struct {
	StationLocationId graphql.ID `json:"stationLocationId"`
	Firstname         string     `json:"firstname"`
	Lastname          string     `json:"lastname"`
	Msisdn            string     `json:"msisdn"`
}

type CreateStationOfficerData struct {
	StationLocationId uuid.UUID `json:"stationLocationId"`
	ProjectId         uuid.UUID `json:"projectId"`
	Firstname         string    `json:"firstname"`
	Lastname          string    `json:"lastname"`
	Msisdn            string    `json:"msisdn"`
	CreatedBy         string    `json:"createdBy"`
}
