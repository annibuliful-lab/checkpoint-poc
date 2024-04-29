package note

import (
	"github.com/graph-gophers/graphql-go"
)

type StationDashboardReportInput struct {
	Note      string     `json:"note"`
}


type StationDashboardReport struct {
	Id        graphql.ID `json:"id"`
	ProjectId graphql.ID `json:"projectId"`
	Note      string     `json:"note"`
}