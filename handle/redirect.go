package handle

import (
  "time"
  "net/http"
  "html/template"

  "github.com/gorilla/mux"
  "github.com/vannio/shrink/db"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("template/index.html")
  token := mux.Vars(r)["token"]
  originalUrl, urlErr := findRow(token)

  if urlErr != nil {
    t.Execute(w, urlErr)
    return
  }

  if len(originalUrl) == 0 {
    http.NotFound(w, r)
    return
  }

  _, queryErr := db.Connection.Exec(
    "UPDATE urls SET access_count = access_count + 1, last_accessed = $1 WHERE token = $2",
    time.Now(),
    token,
  )

  if queryErr != nil {
    t.Execute(w, queryErr)
    return
  }

  http.Redirect(w, r, originalUrl, 302)
}
