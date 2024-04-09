package tag

import (
	"checkpoint/.gen/checkpoint/public/model"

	"github.com/graph-gophers/graphql-go"
)

type UpsertTagData struct {
	Tag       string `json:"tag"`
	ProjectId string `json:"projectId"`
	CreatedBy string `json:"createdBy"`
}

type Tag struct {
	Id        graphql.ID `json:"id"`
	Title     string     `json:"title"`
	ProjectId graphql.ID `json:"projectId"`
}

type GetTagsInput struct {
	Search *string `json:"search"`
	Limit  int64   `json:"limit"`
	Skip   int64   `json:"skip"`
}

type ImeiTag struct {
	model.Tag
	model.ImeiConfigurationTag
}

type ImsiTag struct {
	model.Tag
	model.ImsiConfigurationTag
}

type MobileDeviceTag struct {
	model.Tag
	model.MobileDeviceConfigurationTag
}

type StationLocationTag struct {
	model.StationLocationTag
	model.Tag
}

type VehicleTargetConfigurationTag struct {
	model.VehicleTargetConfigurationTag
	model.Tag
}
