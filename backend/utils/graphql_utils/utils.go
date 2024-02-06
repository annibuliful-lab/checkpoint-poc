package graphql_utils

import "github.com/graph-gophers/graphql-go"

func ConvertStringToNullID(strPtr *string) graphql.NullID {
	if strPtr == nil {
		return graphql.NullID{} // or whatever default value you prefer for NullID
	}
	id := graphql.ID(*strPtr)
	return graphql.NullID{Value: &id}
}
