package url

import (
	"fmt"
	"hash/adler32"
	"os"

	"github.com/PuerkitoBio/purell"
)

// Normalise : Removes trailing slashes and a couple of other things that I don't remember configuring
func Normalise(url string) string {
	return purell.MustNormalizeURLString(
		url,
		purell.FlagsUsuallySafeGreedy|
			purell.FlagRemoveDuplicateSlashes|
			purell.FlagSortQuery,
	)
}

// Make : Puts together a full URL with a given slug
func Make(slug string) string {
	return "http://" + os.Getenv("base_URL") + os.Getenv("port") + "/" + slug
}

// Slug : Hashes a given URL and returns the resulting string as a slug
func Slug(url string) string {
	b := []byte(url)
	c := adler32.Checksum(b)
	return fmt.Sprintf("%x", c)
}
