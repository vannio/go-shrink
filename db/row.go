package db

import (
	"database/sql"
	"time"

	"github.com/vannio/shrink/url"
)

// Row : Data structure for a db row
type Row struct {
	Slug        string
	URL         string
	CreatedAt   time.Time
	AccessCount int
}

// IncrementAccessCount : Tracks short url usage
func (r Row) IncrementAccessCount() error {
	_, err := psql.Exec(
		"UPDATE urls SET access_count = access_count + 1, last_accessed = $1 WHERE slug = $2",
		time.Now(),
		r.Slug,
	)
	return err
}

// FindBySlug : Queries db
func (r Row) FindBySlug(slug string) (Row, error) {
	r.Slug = slug
	err := r.find()
	return r, err
}

// FindByURL : Queries db
func (r Row) FindByURL(normalisedURL string) (Row, error) {
	r.Slug = url.Slug(normalisedURL)
	err := r.find()
	return r, err
}

// Inject : Adds a new entry to the db
func (r Row) Inject(normalisedURL string) error {
	slug := url.Slug(normalisedURL)
	_, err := psql.Exec(
		"INSERT INTO urls(slug,url) VALUES($1,$2) returning id;",
		slug,
		normalisedURL,
	)
	return err
}

func (r *Row) find() error {
	query := "SELECT url, created_at, access_count FROM urls WHERE slug = $1"
	err := psql.QueryRow(query, r.Slug).Scan(&r.URL, &r.CreatedAt, &r.AccessCount)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
