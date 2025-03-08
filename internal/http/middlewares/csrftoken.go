package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

func CSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-CSRF-Token", csrf.Token(c.Request))
	}
}
