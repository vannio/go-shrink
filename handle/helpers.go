package handle

import (
	"database/sql"
	"fmt"
	"hash/adler32"

	"github.com/PuerkitoBio/purell"
	"github.com/vannio/shrink/db"
)

func findRow(slug string) (string, error) {
	var url string
	err := db.Connection.QueryRow("SELECT url FROM urls WHERE slug = $1", slug).Scan(&url)

	if err == sql.ErrNoRows {
		return url, nil
	}

	return url, err
}

func createSlug(url string) string {
	b := []byte(url)
	c := adler32.Checksum(b)
	return fmt.Sprintf("%x", c)
}

func normaliseURL(url string) string {
	return purell.MustNormalizeURLString(
		url,
		purell.FlagsUsuallySafeGreedy|
			purell.FlagRemoveDuplicateSlashes|
			purell.FlagSortQuery,
	)
}
