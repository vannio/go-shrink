package handle

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vannio/shrink/db"
)

// Redirect : This handles the redirection of a shortURL
func Redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	row, err := db.FindRowBySlug(db.Row{}, mux.Vars(r)["slug"])

	if err != nil {
		w.Write(jsonifyError(err))
		return
	}

	if row.URL == "" {
		http.NotFound(w, r)
		return
	}

	err = db.IncrementAccessCount(row)

	if err != nil {
		w.Write(jsonifyError(err))
		return
	}

	http.Redirect(w, r, row.URL, 302)
}
