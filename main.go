package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vannio/shrink/handle"
)

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix("/s").Subrouter().StrictSlash(true)
	s.HandleFunc("/", handle.Root)
	s.HandleFunc("/create", handle.Create)
	s.HandleFunc("/{token}", handle.Redirect)

	http.ListenAndServe(":9000", r)
}
