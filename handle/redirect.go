package handle

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vannio/shrink/db"
)

// Redirect : This handles the redirection of a shortURL
func Redirect(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	originalURL, err := db.FindRow(slug)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.Write(jsonifyError(err))
		return
	}

	if len(originalURL) == 0 {
		http.NotFound(w, r)
		return
	}

	err = db.IncrementVisits(slug)

	if err != nil {
		w.Write(jsonifyError(err))
		return
	}

	http.Redirect(w, r, originalURL, 302)
}
