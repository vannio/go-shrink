package handle

import (
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/vannio/shrink/db"
)

// Create : This handles the creation of a shortURL
func Create(w http.ResponseWriter, r *http.Request) {
	baseURL := os.Getenv("baseURL")
	port := os.Getenv("port")
	pathPrefix := os.Getenv("pathPrefix")

	if r.Method != "POST" {
		http.Redirect(w, r, pathPrefix, 301)
	}

	t, _ := template.ParseFiles("template/index.html")
	u, err := url.ParseRequestURI(r.FormValue("url"))

	if err != nil {
		t.Execute(w, err)
		return
	}

	queryURL := u.String()
	normalisedURL := normaliseURL(queryURL)
	slug := createSlug(normalisedURL)

	if strings.Contains(queryURL, baseURL) && strings.Contains(queryURL, pathPrefix) {
		slug = strings.TrimPrefix(u.EscapedPath(), pathPrefix)
	}

	shortURL := "http://" + baseURL + port + pathPrefix + slug
	originalURL, err := findRow(slug)

	if err != nil {
		t.Execute(w, err)
		return
	}

	if len(originalURL) > 0 {
		if queryURL == shortURL {
			t.Execute(w, "Shorturl already exists! Shorturl for "+originalURL+" is "+shortURL)
			return
		}

		t.Execute(w, "Shorturl already exists! Shorturl for "+queryURL+" is "+shortURL)
		return
	}

	_, err = db.Connection.Exec(
		"INSERT INTO urls(slug,url,created_at) VALUES($1,$2,$3) returning id;",
		slug,
		normalisedURL,
		time.Now(),
	)

	if err != nil {
		t.Execute(w, err)
		return
	}

	t.Execute(w, "Shorturl created! Shorturl for "+queryURL+" is "+shortURL)
}
