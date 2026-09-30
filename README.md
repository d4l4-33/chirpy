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
Clone the repo:
```
git clone https://github.com/primusprag/chirpy
```

You then need to create a .env file in your working folder, it should have the following variables:

DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable *or other database url*
PLATFORM=`dev`
SECRET= a 64 bit random string, you can use following in your terminal:
```
openssl rand -base64 64
```
POLKAKEY= *simulates a webhook API-key, use a 24 bit random string*
```
openssl rand -base64 24
```

## Endpoints
### Admin & Readiness
`GET /api/healthz` - Check server readiness

`GET /admin/metrics` - Check serverhits since reset

`POST /admin/reset` - Resets serverhits and database *(supposed to only work in dev environment)*

### Users
`POST /api/users` - Create new user: requires json in request:
```json
{
    "password": "<password>",
    "email": "<email>"

}
```

`POST /api/login` - Logs in user. Json:
```json
{
    "password": "<password>",
    "email": "<email>"

}
```
Response: *token needed for many of the following endpoints*
```json
{
    "id": ,
    "created_at": ,
    "updated_at": ,
    "email": ,
    "token": ,
    "refresh_token": ,
    "is_chirpy_red": "simulating a premium version of the account"
}
```

`PUT /api/users` - Updates user info if logged in. Json:
```
Headers: Authorization: Bearer <token>

{
    "password": "<password>",
    "email": "<email>"

}
```

### Chirps
`POST /api/chirps` - Creates chirp. Requires Json:
```json
Headers: Authorization: Bearer <token>

{
    "body": "<chirp message>"

}
```

`GET /api/chirps/` alt `GET /api/chirps?author_id={author_id}&sort=desc` - Retrieves chirps by optional {author_id} and the ability to sort time of creation descending, ascending by default.

`GET /api/chirps/{chirpID}` - Retrievs a specific chirp by it's id.

`DELETE /api/chirps/{chirpID}` - Deletes a chirp by it's id.
```
Headers: Authorization: Bearer <token>
```


### Tokens
`POST /api/refresh` - Refeshes the user's API-token. Json

```
Headers: Authorization: Bearer <token>
```

`POST /api/revoke` - Revokes a user's refresh-token requiring them to login anew.
```
Headers: Authorization: Bearer <token>
```

### Webhooks
`POST /api/polka/webhooks` - Updates user's premium status (is_chirpy_red) when recieving the correct request from our "third-party payment-supplier" (*it's a simulation*).
```json
Authorization: APIKey <POLKAKEY in env>

{
    "data": {
        "user_id": "<user_id>"
    },
    "event": "user.upgraded"
}
```

## Ending thoughts
I think that's it for now. Thank you if you read and if not it's been good practice writing this markdown.
And with that, it's a cheers and a cheerio...