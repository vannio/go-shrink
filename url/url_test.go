package url

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNormalisedURLs(t *testing.T) {
	normalisedURL := Normalise("https://www.testing.com///?z=1&a=2")
	if !reflect.DeepEqual(normalisedURL, "https://www.testing.com?a=2&z=1") {
		t.Fail()
	}
}

func TestMakeURL(t *testing.T) {
	url := Make("abc123")

	if !reflect.DeepEqual(url, "https://testing.com/abc123") {
		fmt.Println(url)
		t.Fail()
	}
}

func TestSlugFromURL(t *testing.T) {
	slug := Slug("https://www.testing.com")

	if !reflect.DeepEqual(slug, "694008ca") {
		t.Fail()
	}
}
