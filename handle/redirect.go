package handle

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vannio/shrink/db"
)

// Redirect : This handles the redirection of a shortURL
func Redirect(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/index.html")
	slug := mux.Vars(r)["slug"]
	originalURL, err := db.FindRow(slug)

	if err != nil {
		t.Execute(w, err)
		return
	}

	if len(originalURL) == 0 {
		http.NotFound(w, r)
		return
	}

	err = db.IncrementVisits(slug)

	if err != nil {
		t.Execute(w, err)
		return
	}

	http.Redirect(w, r, originalURL, 302)
}
