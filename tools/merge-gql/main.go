package main

import (
	"os"

	gql "github.com/mununki/gqlmerge/lib"
)

func main() {
	schema := gql.Merge(" ", "./backend")
	err := os.WriteFile("./tools/genql/generated.graphql", []byte(*schema), 0777)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("generated.graphql", []byte(*schema), 0777)
	if err != nil {
		panic(err)
	}
}
