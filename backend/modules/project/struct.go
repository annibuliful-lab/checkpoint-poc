package project

import (
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type Project struct {
	Title     string            `json:"title"`
	ID        graphql.ID        `json:"id"`
	CreatedBy graphql.ID        `json:"createdBy"`
	CreatedAt graphql.Time      `json:"createdAt"`
	UpdatedBy *graphql.NullID   `json:"updatedBy"`
	UpdatedAt *graphql.NullTime `json:"updatedAt"`
}

type CreateProjectInput struct {
	Title     string `json:"title"`
	AccountId string `json:"accountId"`
}

type UpdateProjectInput struct {
	Id        graphql.ID `json:"id"`
	Title     string     `json:"title"`
	AccountId string     `json:"accountId"`
}
type UpdateProjectData struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	AccountId string    `json:"accountId"`
}

type GetProjectInput struct {
	ID graphql.ID `json:"id"`
}

type DeleteProjectInput struct {
	ID        graphql.ID `json:"id"`
	AccountId string     `json:"accountId"`
}
