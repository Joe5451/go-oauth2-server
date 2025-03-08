package http

import (
	"github.com/Joe5451/go-oauth2-server/internal/adapter/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	// setup csrf middleware
	router.Use(middlewares.CSRF())
	router.Use(middlewares.CSRFToken())

	router.GET("/csrf-token", userHandler.CSRFToken)

	return router
}
