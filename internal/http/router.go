package http

import (
	"github.com/Joe5451/go-oauth2-server/internal/adapter/handlers"
	"github.com/Joe5451/go-oauth2-server/internal/http/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	// Set up a cookie-based session temporarily
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("usersession", store))

	// setup csrf middleware
	router.Use(middlewares.CSRF())
	router.Use(middlewares.CSRFToken())

	router.GET("/csrf-token", userHandler.CSRFToken)

	return router
}
