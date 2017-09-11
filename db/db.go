package db

import (
	"database/sql"
	"log"
	"os"

	// Blank import
	_ "github.com/lib/pq"
)

var psql *sql.DB

// AccessDB : Interface to allow mocking DB data
type AccessDB interface {
	FindBySlug(string) (Row, error)
	FindByURL(string) (Row, error)
	Inject(string) error
	IncrementAccessCount() error
}

// FindRowBySlug : Wrapper for FindBySlug
func FindRowBySlug(a AccessDB, slug string) (Row, error) {
	return a.FindBySlug(slug)
}

// FindRowByURL : Wrapper for FindByURL
func FindRowByURL(a AccessDB, url string) (Row, error) {
	return a.FindByURL(url)
}

// InjectRow :Wrapper for InjectRow
func InjectRow(a AccessDB, url string) error {
	return a.Inject(url)
}

// IncrementAccessCount : Wrapper for IncrementAccessCount
func IncrementAccessCount(a AccessDB) error {
	return a.IncrementAccessCount()
}

func init() {
	if os.Getenv("ENV") != "test" {
		var err error
		psql, err = sql.Open("postgres", "dbname=shrink sslmode=disable")

		if err != nil {
			log.Fatal(err)
		}

		if err = psql.Ping(); err != nil {
			log.Fatal(err)
		}
	}
}
