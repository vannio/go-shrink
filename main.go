package main

import (
  "log"
  "time"
  "fmt"
  "net/http"
  "net/url"
  "hash/adler32"
  "database/sql"

  _ "github.com/lib/pq"
  "github.com/PuerkitoBio/purell"
  "github.com/gorilla/mux"
)

var db *sql.DB

func init() {
  var err error
  db, err = sql.Open("postgres", "postgres://localhost/shrink?sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }

  if err = db.Ping(); err != nil {
    log.Fatal(err)
  }
}

func FindOrCreateShorturl(res http.ResponseWriter, req *http.Request) {
  if (req.Method != "POST") {
    http.Redirect(res, req, "/s", 301)
  }

  query := req.FormValue("url")

  if (len(query) == 0) {
    fmt.Fprint(res, "No URL detected")
    return
  }

  _, parseErr := url.ParseRequestURI(query)
  if parseErr != nil {
    fmt.Fprint(res, parseErr)
    return
  }

  normalisedUrl := purell.MustNormalizeURLString(
    query,
    purell.FlagsUsuallySafeGreedy |
    purell.FlagRemoveDuplicateSlashes |
    purell.FlagAddWWW |
    purell.FlagSortQuery,
  )

  token := createToken(normalisedUrl)

  url, urlErr := findRow(token)

  if urlErr != nil {
    fmt.Fprint(res, urlErr)
    return
  }

  fmt.Println(url)

  if len(url) > 0 {
    fmt.Fprint(res, "Shorturl already exists! Shorturl for " + query + " is http://vann.io/s/" + token)
    return
  }

  insertErr := InsertUrlToDB(normalisedUrl, token)

  if insertErr != nil {
    fmt.Fprint(res, insertErr)
    return
  }

  fmt.Fprint(res, "Shorturl created! Shorturl for " + query + " is http://vann.io/s/" + token)
  return
}

func createToken(url string) string {
  b := []byte(url)
  c := adler32.Checksum(b)
  return fmt.Sprintf("%x", c)
}

func InsertUrlToDB(url string, token string) error {
  fmt.Println("# Inserting values")

  var lastInsertId int
  err := db.QueryRow(
    "INSERT INTO urls(token,url,created_at) VALUES($1,$2,$3) returning id;",
    token,
    url,
    time.Now(),
  ).Scan(&lastInsertId)

  if err != nil {
    return err
  }

  fmt.Println("last inserted id =", lastInsertId)
  return nil
}

func findRow(token string) (string, error) {
  var url string
  err := db.QueryRow("SELECT url FROM urls WHERE token = $1", token).Scan(&url)

  if err == sql.ErrNoRows {
    return url, nil
  }

  return url, err
}

func RedirectToUrl(res http.ResponseWriter, req *http.Request) {
  token := mux.Vars(req)["token"]

  url, urlErr := findRow(token)

  if urlErr != nil {
    fmt.Fprint(res, urlErr)
    return
  }

  if len(url) > 0 {
    http.Redirect(res, req, url, 301)
  }

  http.NotFound(res, req)
}

func Homepage(res http.ResponseWriter, req *http.Request) {
  fmt.Println(req.URL.Path[1:])
  http.ServeFile(res, req, "views/index.html")
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/s", Homepage)
  r.HandleFunc("/s/create", FindOrCreateShorturl)
  r.HandleFunc("/s/{token}", RedirectToUrl)

  http.ListenAndServe(":9000", r)
}
