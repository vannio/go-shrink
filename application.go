package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vannio/shrink/handle"
)

func main() {
	baseURL := "localhost"
	pathPrefix := "/s/"
	port := ":8080"

	os.Setenv("port", port)
	os.Setenv("baseURL", baseURL)
	os.Setenv("pathPrefix", pathPrefix)

	r := mux.NewRouter()
	s := r.PathPrefix(pathPrefix).Subrouter().StrictSlash(true)
	s.HandleFunc("/", handle.Root)
	s.HandleFunc("/create", handle.Create)
	s.HandleFunc("/{token}", handle.Redirect)

	fmt.Println("Server listening at", baseURL+port)
	http.ListenAndServe(port, r)
}
