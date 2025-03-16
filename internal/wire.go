//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/Joe5451/go-oauth2-server/internal/adapter/handlers"
	"github.com/Joe5451/go-oauth2-server/internal/adapter/repositories"
	"github.com/Joe5451/go-oauth2-server/internal/application"
	"github.com/Joe5451/go-oauth2-server/internal/application/ports/in"
	"github.com/Joe5451/go-oauth2-server/internal/application/ports/out"
	"github.com/Joe5451/go-oauth2-server/internal/database"
	"github.com/Joe5451/go-oauth2-server/internal/http"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var providerSet wire.ProviderSet = wire.NewSet(
	database.NewPostgresDB,

	wire.Bind(new(out.UserRepository), new(*repositories.PostgresUserRepository)),
	repositories.NewPostgresUserRepository,

	wire.Bind(new(in.UserUsecase), new(*application.UserService)),
	application.NewUserService,

	handlers.NewUserHandler,
	handlers.NewTemplateHandler,

	http.NewRouter,
)

func InitializeApp() (*gin.Engine, error) {
	panic(
		wire.Build(
			providerSet,
		),
	)
}
