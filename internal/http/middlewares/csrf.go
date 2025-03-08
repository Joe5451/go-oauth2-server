package middlewares

import (
	"net/http"

	"github.com/Joe5451/go-oauth2-server/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

func CSRF() gin.HandlerFunc {
	csrfMiddleware := csrf.Protect(
		[]byte("32-byte-long-auth-key"),
		csrf.Path("/"),
		csrf.HttpOnly(true),
		csrf.Secure(config.AppConfig.CSRFSecure),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message": "Forbidden - CSRF token invalid"}`))
		})),
	)

	return adapter.Wrap(csrfMiddleware)
}
