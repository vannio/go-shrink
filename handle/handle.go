package handle

import (
	"encoding/json"

	"github.com/vannio/shrink/db"
	"github.com/vannio/shrink/url"
)

// Response : Structure of the JSON
type Response struct {
	Error    string `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
	QueryURL string `json:"query,omitempty"`
	ShortURL string `json:"shorturl,omitempty"`
}

func shrink(queryURL string) Response {
	normalisedURL := url.Normalise(queryURL)
	slug := url.Slug(normalisedURL)
	shortURL := url.Make(slug)
	originalURL, err := db.FindRow(slug)

	if err != nil {
		return Response{Error: err.Error()}
	}

	if len(originalURL) > 0 {
		url := queryURL

		if queryURL == shortURL {
			url = originalURL
		}

		return Response{
			Message:  "Already exists",
			QueryURL: url,
			ShortURL: shortURL,
		}
	}

	err = db.AddRow(slug, normalisedURL)

	if err != nil {
		return Response{Error: err.Error()}
	}

	return Response{
		Message:  "Shorturl created",
		QueryURL: queryURL,
		ShortURL: shortURL,
	}
}

func jsonifyError(err error) []byte {
	response := Response{Error: err.Error()}
	jsonData, _ := json.Marshal(response)
	return jsonData
}
