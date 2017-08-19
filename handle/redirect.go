package handle

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vannio/shrink/db"
)

// Redirect : This handles the redirection of a shortURL
func Redirect(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/index.html")
	slug := mux.Vars(r)["slug"]
	originalURL, urlErr := findRow(slug)

	if urlErr != nil {
		t.Execute(w, urlErr)
		return
	}

	if len(originalURL) == 0 {
		http.NotFound(w, r)
		return
	}

	_, queryErr := db.Connection.Exec(
		"UPDATE urls SET access_count = access_count + 1, last_accessed = $1 WHERE slug = $2",
		time.Now(),
		slug,
	)

	if queryErr != nil {
		t.Execute(w, queryErr)
		return
	}

	http.Redirect(w, r, originalURL, 302)
}
