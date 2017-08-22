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
	port := ":8080"

	os.Setenv("port", port)
	os.Setenv("base_URL", baseURL)

	r := mux.NewRouter()
	r.HandleFunc("/", handle.Root)
	r.HandleFunc("/create", handle.Create)
	r.HandleFunc("/{slug}", handle.Redirect)

	fmt.Println("Server listening at", baseURL+port)
	http.ListenAndServe(port, r)
}
