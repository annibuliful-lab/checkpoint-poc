package db

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	_db  *sql.DB
	once sync.Once
)

func GetClient() *sql.DB {
	once.Do(func() {

		connectString := os.Getenv("PRIMARY_DATABASE_URL")
		database, err := sql.Open("postgres", connectString)

		if err != nil {
			log.Fatal(err)
		}
		database.SetMaxOpenConns(5)
		database.SetMaxIdleConns(5)
		_db = database

	})

	return _db
}
