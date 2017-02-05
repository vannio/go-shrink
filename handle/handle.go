package handle

import (
  "time"
  "fmt"
  "net/http"
  "net/url"
  "hash/adler32"
  "database/sql"
  "html/template"

  _ "github.com/lib/pq"
  "github.com/gorilla/mux"
  "github.com/PuerkitoBio/purell"
  "github.com/vannio/shrink/db"
)

// ------
// PUBLIC
// ------

func Create(w http.ResponseWriter, r *http.Request) {
  if (r.Method != "POST") {
    http.Redirect(w, r, "/s", 301)
  }

  query := r.FormValue("url")
  t, _ := template.ParseFiles("template/index.html")

  // URL validation
  _, parseErr := url.ParseRequestURI(query)

  if parseErr != nil {
    t.Execute(w, parseErr)
    return
  }

  // Sanitize URLs to avoid simple duplicates
  normalisedUrl := normaliseUrl(query)

  // Generate (hopefully) unique token
  token := createToken(normalisedUrl)

  // Check if URL has already been submitted
  originalUrl, urlErr := findRow(token)

  if urlErr != nil {
    t.Execute(w, urlErr)
    return
  }

  if len(originalUrl) > 0 {
    t.Execute(w, "Shorturl already exists! Shorturl for " + query + " is http://vann.io/s/" + token)
    return
  }

  // Create new entry in db
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

func Redirect(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("template/index.html")
  token := mux.Vars(r)["token"]
  originalUrl, urlErr := findRow(token)

  if urlErr != nil {
    t.Execute(w, urlErr)
    return
  }

  // 404 when token is invalid
  if len(originalUrl) == 0 {
    http.NotFound(w, r)
    return
  }

  // Otherwise update access_count and last_accessed
  _, queryErr := db.Connection.Exec(
    "UPDATE urls SET access_count = access_count + 1, last_accessed = $1 WHERE token = $2",
    time.Now(),
    token,
  )

  if queryErr != nil {
    t.Execute(w, queryErr)
    return
  }

  // Redirect to the original URL
  http.Redirect(w, r, originalUrl, 302)
}

func Root(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("template/index.html")
  t.Execute(w, nil)
}

// -------
// PRIVATE
// -------

func findRow(token string) (string, error) {
  var url string
  err := db.Connection.QueryRow("SELECT url FROM urls WHERE token = $1", token).Scan(&url)

  if err == sql.ErrNoRows {
    return url, nil
  }

  return url, err
}

func createToken(url string) string {
  b := []byte(url)
  c := adler32.Checksum(b)
  return fmt.Sprintf("%x", c)
}

func normaliseUrl(url string) string {
  return purell.MustNormalizeURLString(
    url,
    purell.FlagsUsuallySafeGreedy |
    purell.FlagRemoveDuplicateSlashes |
    purell.FlagAddWWW |
    purell.FlagSortQuery,
  )
}
