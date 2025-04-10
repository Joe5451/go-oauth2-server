# Go Oauth2 Server

[![Go Report Card](https://goreportcard.com/badge/github.com/Joe5451/go-oauth2-server)](https://goreportcard.com/report/github.com/Joe5451/go-oauth2-server)

## Project Layout

```
.
├── api
├── assets
├── build
├── cmd
├── deployments
├── internal/
│   ├── adapter/
│   │   ├── handlers/
│   │   │   └── user_handler.go
│   │   └── repositories/
│   │       └── postgres_user_repository.go
│   ├── application/
│   │   ├── user_service.go
│   │   └── ports/
│   │       ├── in/
│   │       │   └── user_usecase.go
│   │       └── out/
│   │           └── user_repository.go
│   ├── config
│   ├── constants
│   ├── domain/
│   │   ├── user.go
│   │   └── social_account.go
│   ├── http/
│   │   ├── middlewares
│   │   └── router.go
│   ├── wire_gen.go
│   └── wire.go
├── test
├── go.mod
├── go.sum
└── README.md
```

## Test

Copy `.env.example` to `.env.test` and modify `.env.test` to match your testing environment settings.
```bash
cp .env.example .env.test
```

**Running Integration Tests**
```
go test ./test
```
