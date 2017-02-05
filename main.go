package main

import (
  "log"
  "time"
  "fmt"
  "net/http"
  "net/url"
  "hash/adler32"
  "database/sql"
  "html/template"

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
  t, _ := template.ParseFiles("views/index.html")

  if (len(query) == 0) {
    t.Execute(res, "Please submit a valid URL")
    return
  }

  _, parseErr := url.ParseRequestURI(query)
  if parseErr != nil {
    t.Execute(res, parseErr)
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
    t.Execute(res, urlErr)
    return
  }

  if len(url) > 0 {
    t.Execute(res, "Shorturl already exists! Shorturl for " + query + " is http://vann.io/s/" + token)
    return
  }

  insertErr := InsertUrlToDB(normalisedUrl, token)

  if insertErr != nil {
    t.Execute(res, insertErr)
    return
  }

  t.Execute(res, "Shorturl created! Shorturl for " + query + " is http://vann.io/s/" + token)
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
  t, _ := template.ParseFiles("views/index.html")

  if urlErr != nil {
    t.Execute(res, urlErr)
    return
  }

  if len(url) > 0 {
    http.Redirect(res, req, url, 301)
  }

  http.NotFound(res, req)
}

func Homepage(res http.ResponseWriter, req *http.Request) {
  t, _ := template.ParseFiles("views/index.html")
  t.Execute(res, nil)
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/s", Homepage)
  r.HandleFunc("/s/create", FindOrCreateShorturl)
  r.HandleFunc("/s/{token}", RedirectToUrl)

  http.ListenAndServe(":9000", r)
}
