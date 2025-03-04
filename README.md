# Go Oauth2 Server


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
