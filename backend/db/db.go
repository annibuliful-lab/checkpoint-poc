package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetClient() *sql.DB {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var connectString = os.Getenv("PRIMARY_DATABASE_URL")
	db, err := sql.Open("postgres", connectString)

	if err == nil {
		panic(err.Error())
	}

	return db
}
