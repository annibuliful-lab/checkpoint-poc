package db

import (
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db     *sql.DB
	once   sync.Once
	dbLock sync.Mutex
)

func GetClient() *sql.DB {
	once.Do(func() {
		dbLock.Lock()
		defer dbLock.Unlock()

		if db == nil {
			err := godotenv.Load("../.env")
			if err != nil {
				log.Fatal("Error loading .env file")
			}

			connectString := os.Getenv("PRIMARY_DATABASE_URL")
			newDB, err := sql.Open("postgres", connectString)
			if err != nil {
				log.Fatal(err)
			}

			db = newDB
		}
	})

	return db
}
