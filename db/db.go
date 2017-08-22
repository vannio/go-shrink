package db

import (
	"database/sql"
	"log"
	"time"

	// Blank import
	_ "github.com/lib/pq"
)

// Connection : Exposes the open database connection globally
var Connection *sql.DB

// FindRow : Queries db to see if slug already exists
func FindRow(slug string) (string, error) {
	var url string
	err := Connection.QueryRow("SELECT url FROM urls WHERE slug = $1", slug).Scan(&url)

	if err == sql.ErrNoRows {
		return url, nil
	}

	return url, err
}

// AddRow : Adds a new entry to the db
func AddRow(slug string, url string) error {
	_, err := Connection.Exec(
		"INSERT INTO urls(slug,url,created_at) VALUES($1,$2,$3) returning id;",
		slug,
		url,
		time.Now(),
	)

	return err
}

// IncrementVisits : Tracks short url usage
func IncrementVisits(slug string) error {
	_, err := Connection.Exec(
		"UPDATE urls SET access_count = access_count + 1, last_accessed = $1 WHERE slug = $2",
		time.Now(),
		slug,
	)

	return err
}

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
