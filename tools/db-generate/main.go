package main

import (
	"github.com/go-jet/jet/v2/generator/postgres"
	_ "github.com/lib/pq"
)

func main() {

	err := postgres.Generate("../../backend/.gen/",
		postgres.DBConnection{
			Host:       "localhost",
			Port:       5432,
			User:       "checkpoint",
			Password:   "checkpoint",
			DBName:     "checkpoint",
			SchemaName: "public",
			SslMode:    "disable",
		})

	if err != nil {
		panic(err)
	}
}
