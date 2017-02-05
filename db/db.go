package db

import (
  "log"
  "database/sql"
)

// DB connection is safe to expose
var Connection *sql.DB

func init() {
  var err error
  Connection, err = sql.Open("postgres", "postgres://localhost/shrink?sslmode=disable")

  if err != nil {
    log.Fatal(err)
  }

  if err = Connection.Ping(); err != nil {
    log.Fatal(err)
  }
}
