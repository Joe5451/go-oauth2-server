package http

import (
	"github.com/Joe5451/go-oauth2-server/internal/adapter/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	return router
}
