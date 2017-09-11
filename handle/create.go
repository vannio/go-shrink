package handle

import (
	"encoding/json"
	"net/http"
	nurl "net/url"

	"github.com/vannio/shrink/db"
	"github.com/vannio/shrink/url"
)

// Create : This handles the creation of a shortURL
func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	u, err := nurl.ParseRequestURI(r.FormValue("url"))
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.Write(jsonifyError(err))
		return
	}

	queryURL := u.String()
	response := shrink(queryURL)
	jsonData, _ := json.Marshal(response)
	w.Write(jsonData)
}

func shrink(queryURL string) Response {
	normalisedURL := url.Normalise(queryURL)
	row, err := db.FindRowByURL(db.Row{}, normalisedURL)

	if err != nil {
		return Response{Error: err.Error()}
	}

	response := Response{
		QueryURL: queryURL,
		ShortURL: url.Make(row.Slug),
	}

	if row.URL != "" {
		response.SetMessage("Already exists")
		return response
	}

	err = db.InjectRow(row, normalisedURL)

	if err != nil {
		return Response{Error: err.Error()}
	}

	response.SetMessage("Shorturl created")
	return response
}
