package handle

import (
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/vannio/shrink/db"
)

// Create : This handles the creation of a shortURL
func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/s", 301)
	}

	query := r.FormValue("url")
	t, _ := template.ParseFiles("template/index.html")

	_, parseErr := url.ParseRequestURI(query)

	if parseErr != nil {
		t.Execute(w, parseErr)
		return
	}

	normalisedURL := normaliseURL(query)

	token := createToken(normalisedURL)

	originalURL, urlErr := findRow(token)

	if urlErr != nil {
		t.Execute(w, urlErr)
		return
	}

	if len(originalURL) > 0 {
		t.Execute(w, "Shorturl already exists! Shorturl for "+query+" is http://vann.io/s/"+token)
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

	t.Execute(w, "Shorturl created! Shorturl for "+query+" is http://vann.io/s/"+token)
}
