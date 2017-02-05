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

func Create(res http.ResponseWriter, req *http.Request) {
  if (req.Method != "POST") {
    http.Redirect(res, req, "/s", 301)
  }

  query := req.FormValue("url")
  t, _ := template.ParseFiles("tmpl/index.html")

  // URL validation
  _, parseErr := url.ParseRequestURI(query)

  if parseErr != nil {
    t.Execute(res, parseErr)
    return
  }

  // Sanitize URLs to avoid simple duplicates
  normalisedUrl := normaliseUrl(query)

  // Generate (hopefully) unique token
  token := createToken(normalisedUrl)

  // Check if URL has already been submitted
  originalUrl, urlErr := findRow(token)

  if urlErr != nil {
    t.Execute(res, urlErr)
    return
  }

  if len(originalUrl) > 0 {
    t.Execute(res, "Shorturl already exists! Shorturl for " + query + " is http://vann.io/s/" + token)
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
    t.Execute(res, insertErr)
    return
  }

  t.Execute(res, "Shorturl created! Shorturl for " + query + " is http://vann.io/s/" + token)
}

func Redirect(res http.ResponseWriter, req *http.Request) {
  t, _ := template.ParseFiles("tmpl/index.html")
  token := mux.Vars(req)["token"]
  originalUrl, urlErr := findRow(token)

  if urlErr != nil {
    t.Execute(res, urlErr)
    return
  }

  // 404 when token is invalid
  if len(originalUrl) == 0 {
    http.NotFound(res, req)
    return
  }

  // Otherwise update access_count and last_accessed
  _, queryErr := db.Connection.Exec(
    "UPDATE urls SET access_count = access_count + 1, last_accessed = $1 WHERE token = $2",
    time.Now(),
    token,
  )

  if queryErr != nil {
    t.Execute(res, queryErr)
    return
  }

  // Redirect to the original URL
  http.Redirect(res, req, originalUrl, 302)
}

func Root(res http.ResponseWriter, req *http.Request) {
  t, _ := template.ParseFiles("tmpl/index.html")
  t.Execute(res, nil)
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
