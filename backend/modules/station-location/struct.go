package stationlocation

import (
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type StationLocation struct {
	ID          graphql.ID `json:"id"`
	ProjectId   graphql.ID `json:"projectId"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Department  string     `json:"department"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	Remark      *string    `json:"remark"`
}

type DeleteStationLocationInput struct {
	ID graphql.ID `json:"id"`
}

type DeleteStationLocationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	DeletedBy string    `json:"deletedBy"`
}

type GetStationLocationsData struct {
	ProjectId uuid.UUID `json:"projectId"`
	Limit     int64     `json:"limit"`
	Skip      int64     `json:"skip"`
	Search    *string   `json:"search"`
	Tags      *[]string `json:"tags"`
}

type GetStationLocationsInput struct {
	Limit  float64   `json:"limit"`
	Skip   float64   `json:"skip"`
	Search *string   `json:"search"`
	Tags   *[]string `json:"tags"`
}

type GetStationLocationByIdInput struct {
	Id graphql.ID `json:"id"`
}

type GetStationLocationByIdData struct {
	Id        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
}

type UpdateStationLocationInput struct {
	Id          graphql.ID                  `json:"id"`
	Title       *string                     `json:"title"`
	Description *string                     `json:"description"`
	Department  *string                     `json:"department"`
	Latitude    *float64                    `json:"latitude"`
	Longitude   *float64                    `json:"longitude"`
	Remark      *string                     `json:"remark"`
	Tags        *[]string                   `json:"tags"`
	Officers    *[]UpserStationOfficerInput `json:"officerInputs"`
}

type UpdateStationLocationData struct {
	Id          uuid.UUID `json:"id"`
	ProjectId   uuid.UUID `json:"projectId"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Department  *string   `json:"department"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	Remark      *string   `json:"remark"`
	Tags        *[]string `json:"tags"`
	UpdatedBy   string    `json:"updatedBy"`
}

type UpserStationOfficerInput struct {
	Firstname string
	Lastname  *string
	Msisdn    string
}

type CreateStationLocationInput struct {
	Title       string                      `json:"title"`
	Description *string                     `json:"description"`
	Department  string                      `json:"department"`
	Latitude    float64                     `json:"latitude"`
	Longitude   float64                     `json:"longitude"`
	Remark      *string                     `json:"remark"`
	Tags        *[]string                   `json:"tags"`
	Officers    *[]UpserStationOfficerInput `json:"officerInputs"`
}

type CreateStationLocationData struct {
	ProjectId   string    `json:"projectId"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Department  string    `json:"department"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Remark      *string   `json:"remark"`
	CreatedBy   string    `json:"createdBy"`
	Tags        *[]string `json:"tags"`
}
