package handle

import (
	"encoding/json"
	"net/http"
	nurl "net/url"
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
