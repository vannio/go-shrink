package handle

import (
	"html/template"
	"net/http"
)

// Root : This handles the homepage (root)
func Root(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/index.html")
	t.Execute(w, nil)
}
