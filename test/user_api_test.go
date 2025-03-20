package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Joe5451/go-oauth2-server/internal"
	"github.com/Joe5451/go-oauth2-server/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	router *gin.Engine
)

func setup() {
	var err error
	err = os.Chdir("../")
	if err != nil {
		panic(err)
	}

	viper.Set("DB_HOST", "localhost")
	viper.Set("DB_PORT", "5432")
	viper.Set("DB_USER", "postgres")
	viper.Set("DB_PASSWORD", "postgres")
	viper.Set("DB_NAME", "go-oauth2-server")

	viper.Set("REDIS_HOST", "localhost")
	viper.Set("REDIS_PORT", "6379")
	viper.Set("REDIS_PASSWORD", "")
	viper.Set("REDIS_SECRET", "")

	viper.Set("GOOGLE_OAUTH2_CLIENT_ID", "test_google_client_id")
	viper.Set("GOOGLE_OAUTH2_CLIENT_SECRET", "test_google_client_secret")

	viper.Set("FACEBOOK_OAUTH2_CLIENT_ID", "test_facebook_client_id")
	viper.Set("FACEBOOK_OAUTH2_CLIENT_SECRET", "test_facebook_client_secret")

	viper.Set("JWT_SECRET_KEY", "jwt-secret")

	viper.Set("CSRF_SECRET_KEY", "csrf-key")
	viper.Set("CSRF_SECURE", "false")

	viper.Set("UPLOAD_BASE_URL", "http://localhost:8000")

	err = viper.Unmarshal(&config.AppConfig)
	if err != nil {
		panic(err)
	}

	router, err = internal.InitializeApp()
	if err != nil {
		panic(err)
	}
}

func TestCSRFToken(t *testing.T) {
	setup()
	req, _ := http.NewRequest("GET", "/csrf-token", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedStatus := http.StatusNoContent
	assert.Equal(t, expectedStatus, w.Code)

	csrfToken := w.Header().Get("X-CSRF-Token")
	assert.NotEmpty(t, csrfToken, "Expected CSRF token to be present in header")
}
