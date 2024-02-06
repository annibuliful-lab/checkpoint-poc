package gql

import gql "github.com/mununki/gqlmerge/lib"

var MergedSchema = gql.Merge(" ", "./backend")
