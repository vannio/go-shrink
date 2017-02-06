package handle

import (
  "time"
  "net/http"
  "net/url"
  "html/template"

  "github.com/vannio/shrink/db"
)

func Create(w http.ResponseWriter, r *http.Request) {
  if (r.Method != "POST") {
    http.Redirect(w, r, "/s", 301)
  }

  query := r.FormValue("url")
  t, _ := template.ParseFiles("template/index.html")

  _, parseErr := url.ParseRequestURI(query)

  if parseErr != nil {
    t.Execute(w, parseErr)
    return
  }

  normalisedUrl := normaliseUrl(query)

  token := createToken(normalisedUrl)

  originalUrl, urlErr := findRow(token)

  if urlErr != nil {
    t.Execute(w, urlErr)
    return
  }

  if len(originalUrl) > 0 {
    t.Execute(w, "Shorturl already exists! Shorturl for " + query + " is http://vann.io/s/" + token)
    return
  }

  _, insertErr := db.Connection.Exec(
    "INSERT INTO urls(token,url,created_at) VALUES($1,$2,$3) returning id;",
    token,
    normalisedUrl,
    time.Now(),
  )

  if insertErr != nil {
    t.Execute(w, insertErr)
    return
  }

  t.Execute(w, "Shorturl created! Shorturl for " + query + " is http://vann.io/s/" + token)
}
