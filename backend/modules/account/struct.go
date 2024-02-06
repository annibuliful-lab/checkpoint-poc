package account

import "github.com/graph-gophers/graphql-go"

type Account struct {
	Id        graphql.ID   `json:"id"`
	Username  string       `json:"username"`
	CreatedAt graphql.Time `json:"createdAt"`
}
