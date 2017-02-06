package handle

import (
  "net/http"
  "html/template"
)

func Root(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("template/index.html")
  t.Execute(w, nil)
}
