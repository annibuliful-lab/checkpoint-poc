package db

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	_db    *sql.DB
	onceDb sync.Once
)

func GetPrimaryClient() *sql.DB {
	onceDb.Do(func() {

		connectString := os.Getenv("PRIMARY_DATABASE_URL")
		database, err := sql.Open("postgres", connectString)

		if err != nil {
			log.Fatal(err)
		}
		database.SetMaxOpenConns(10)
		database.SetMaxIdleConns(5)
		_db = database

	})

	return _db
}
