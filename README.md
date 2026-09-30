# Chirpy README WIP
This is a guided bootdev project for the [Learn HTTP Servers in Go](https://www.boot.dev/courses/learn-http-servers-golang) course.

## Tools
In this course I have used - *among others* - the [net/http](https://pkg.go.dev/net/http), [encoding/json](https://pkg.go.dev/encoding/json), [argon2id](https://github.com/alexedwards/argon2id), [jwt]("github.com/golang-jwt/jwt/v5") and [godotenv]("github.com/joho/godotenv") packages.
I have used [PostgreSQL](https://www.postgresql.org/), [goose](https://github.com/pressly/goose) and [sqlc](https://sqlc.dev/) for handling the database code.

## Summary
The project simulates an HTTP Server by running the server and database on the local machine. The authentication is done by [argon2id](https://github.com/alexedwards/argon2id) using a secret key in an .env file.

## Installation
Install the 


## Quickstart
You need to create a .env file in your working folder, it should have the following variables:

DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable *or other database url*
PLATFORM=dev
SECRET= *a 64 bit random string* you can use following in your terminal:
```
openssl rand -base64 64
```
POLKAKEY= *simulates a webhook API-key, user a base 24 bit random string*
```
openssl rand -base64 24
```

## Commands
Use `go run chirpy` to start the 
