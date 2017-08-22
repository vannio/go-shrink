package handle

import (
	"html/template"
	"net/http"
	nurl "net/url"

	"github.com/vannio/shrink/db"
	"github.com/vannio/shrink/url"
)

func shrink(queryURL string) (string, error) {
	normalisedURL := url.Normalise(queryURL)
	slug := url.Slug(normalisedURL)
	shortURL := url.Make(slug)
	originalURL, err := db.FindRow(slug)

	if err != nil {
		return "", err
	}

	if len(originalURL) > 0 {
		url := queryURL

		if queryURL == shortURL {
			url = originalURL
		}

		return "Shorturl already exists! Shorturl for " + url + " is " + shortURL, nil
	}

	err = db.AddRow(slug, normalisedURL)

	if err != nil {
		return "", err
	}

	return "Shorturl created! Shorturl for " + queryURL + " is " + shortURL, nil
}

// Create : This handles the creation of a shortURL
func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	t, _ := template.ParseFiles("template/index.html")
	u, err := nurl.ParseRequestURI(r.FormValue("url"))

	if err != nil {
		t.Execute(w, err)
		return
	}

	queryURL := u.String()
	msg, err := shrink(queryURL)

	if err != nil {
		t.Execute(w, err)
		return
	}

	t.Execute(w, msg)
}
