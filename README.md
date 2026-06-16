# Ainyx Users API

RESTful API in GoFiber for managing users with `name` and `dob`. The API stores date of birth and calculates age dynamically with Go's `time` package when users are fetched.

## Submission Links

- GitHub repository: https://github.com/Karthisgowda/Ainyx
- Live deployment: https://ainyx-blond.vercel.app

## Tech Stack

- GoFiber
- PostgreSQL
- SQLC
- Uber Zap
- go-playground/validator

## Endpoints

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/users` | Create a user |
| `GET` | `/users/:id` | Get a user with calculated age |
| `PUT` | `/users/:id` | Update a user |
| `DELETE` | `/users/:id` | Delete a user |
| `GET` | `/users?limit=20&offset=0` | List users with calculated ages |
| `GET` | `/` | API status and endpoint list |
| `GET` | `/health` | Health check |

## Local Setup

1. Copy environment values:

```bash
cp .env.example .env
```

2. Start PostgreSQL and the API with Docker:

```bash
docker compose up --build
```

3. API runs on:

```text
http://localhost:8080
```

## Run Without Docker

Create the `users` table using `db/migrations/001_create_users.sql`, then run:

```bash
go mod download
go run ./cmd/server
```

## SQLC

The SQLC config and query files are included:

```bash
sqlc generate
```

Generated-style files are committed under `db/sqlc` so the project can be reviewed even before regenerating locally.

## Example Requests

Create:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```

Get by ID:

```bash
curl http://localhost:8080/users/1
```

Update:

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'
```

Delete:

```bash
curl -X DELETE http://localhost:8080/users/1
```

## Tests

```bash
go test ./...
```

## Deployment Requirements

For a live deployment, provide:

- A deployed PostgreSQL database URL, for example Neon, Supabase, Render Postgres, or Railway Postgres.
- A hosting account that supports Docker or Go services, for example Render, Railway, Fly.io, or Google Cloud Run.
- Environment variables: `APP_ENV=production`, `PORT=8080`, and `DATABASE_URL=<your postgres url>`.
- Run the SQL migration before starting the service.

## Notes

- `dob` is stored as a database `DATE`.
- `age` is calculated dynamically and is not stored.
- Responses include `X-Request-Id`.
- Request duration is logged with Zap.
- `/users` supports pagination through `limit` and `offset`.
