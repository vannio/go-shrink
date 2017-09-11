# Shrink [![Build Status](https://api.travis-ci.org/vannio/go-shrink.svg?branch=master)](https://travis-ci.org/vannio/go-shrink)

### Useful commands

- PSQL —
  ```sql
  CREATE DATABASE shrink;
  CREATE TABLE urls (
      id serial primary key,
      slug text unique not null,
      url text not null,
      created_at timestamp default current_timestamp(2),
      last_accessed timestamp,
      access_count integer default 0
  );
  ```
- Run tests — `ENV=test BASEURL=https://testing.com go test ./...`
- Start app — `go run main.go` _(default host is http://localhost:8080)_
- Shrink! — `curl -X POST -F url=http://www.reallylongaddress.com http://localhost:8080/create`
