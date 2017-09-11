package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vannio/shrink/handle"
)

func main() {
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", ":8080")
	}

	if os.Getenv("BASEURL") == "" {
		os.Setenv("BASEURL", "http://localhost")
	}

	port := os.Getenv("PORT")
	baseURL := os.Getenv("BASEURL")

	r := mux.NewRouter()
	r.HandleFunc("/create", handle.Create)
	r.HandleFunc("/{slug}", handle.Redirect)

	fmt.Println("Server listening at", baseURL+port)
	http.ListenAndServe(port, r)
}
