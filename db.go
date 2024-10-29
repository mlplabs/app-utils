package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

func InitDb(hostName, dbPort, dbName, dbUser, dbPass string) {
	var err error

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", hostName, dbPort, dbName, dbUser, dbPass)
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}

}

func Close() {
	Db.Close()
}
