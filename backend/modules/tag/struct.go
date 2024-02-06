package tag

import "github.com/graph-gophers/graphql-go"

type Tag struct {
	Id        graphql.ID `json:"id"`
	Title     string     `json:"title"`
	ProjectId string     `json:"projectId"`
}
