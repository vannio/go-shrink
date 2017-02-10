package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// Connection : Exposes the open database connection globally
var Connection *sql.DB

func init() {
	var err error
	Connection, err = sql.Open("postgres", "postgres://localhost/shrink?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	if err = Connection.Ping(); err != nil {
		log.Fatal(err)
	}
}
