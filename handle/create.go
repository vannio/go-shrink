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
	u, parseErr := url.ParseRequestURI(r.FormValue("url"))

	if parseErr != nil {
		t.Execute(w, parseErr)
		return
	}

	queryURL := u.String()
	normalisedURL := normaliseURL(queryURL)
	token := createToken(normalisedURL)

	if strings.Contains(queryURL, baseURL) && strings.Contains(queryURL, pathPrefix) {
		token = strings.TrimPrefix(u.EscapedPath(), pathPrefix)
	}

	shortURL := "http://" + baseURL + port + pathPrefix + token
	originalURL, urlErr := findRow(token)

	if urlErr != nil {
		t.Execute(w, urlErr)
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

	_, insertErr := db.Connection.Exec(
		"INSERT INTO urls(token,url,created_at) VALUES($1,$2,$3) returning id;",
		token,
		normalisedURL,
		time.Now(),
	)

	if insertErr != nil {
		t.Execute(w, insertErr)
		return
	}

	t.Execute(w, "Shorturl created! Shorturl for "+queryURL+" is "+shortURL)
}
