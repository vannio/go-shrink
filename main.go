package main

import (
  "log"
  "time"
  "fmt"
  "database/sql"
  _ "github.com/lib/pq"
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

func main() {
  AddUrl("https://google.co.uk")
}

func AddUrl(url string) {
  fmt.Println("# Inserting values")

  var lastInsertId int
  err := db.QueryRow("INSERT INTO urls(token,url,created_at) VALUES($1,$2,$3) returning token;", "j4s56q", url, time.Now()).Scan(&lastInsertId)

  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("last inserted id =", lastInsertId)
}
