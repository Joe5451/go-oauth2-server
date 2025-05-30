package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Joe5451/go-oauth2-server/internal"
	"github.com/Joe5451/go-oauth2-server/internal/config"
	"github.com/Joe5451/go-oauth2-server/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type TestSuite struct {
	suite.Suite
	router    *gin.Engine
	csrfToken string
	cookies   []*http.Cookie
	conn      *pgx.Conn
}

func (s *TestSuite) SetupSuite() {
	s.Require().NoError(os.Chdir("../"))

	viper.SetConfigName(".env.test")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	s.Require().NoError(viper.ReadInConfig(), "Error reading .env.test file")
	s.Require().NoError(viper.Unmarshal(&config.AppConfig), "Error unmarshalling config")

	var err error
	s.router, err = internal.InitializeApp()
	s.Require().NoError(err)

	s.conn, err = database.NewPostgresDB()
	s.Require().NoError(err, "Failed to connect database for cleanup")
}

func (s *TestSuite) SetupTest() {
	req, _ := http.NewRequest("GET", "/api/csrf-token", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Require().Equal(http.StatusNoContent, w.Code)
	s.csrfToken = w.Header().Get("X-CSRF-Token")
	s.Require().NotEmpty(s.csrfToken)
	s.cookies = w.Result().Cookies()
}

func (s *TestSuite) TearDownTest() {
	tx, err := s.conn.Begin(context.Background())
	s.Require().NoError(err, "Failed to start transaction for cleanup")

	_, err = tx.Exec(context.Background(), `
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`)
	s.Require().NoError(err, "Failed to clean database")

	err = tx.Commit(context.Background())
	s.Require().NoError(err, "Failed to commit cleanup transaction")
}

func (s *TestSuite) createTestUser(name, email, password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	s.Require().NoError(err, "Failed to hash password")

	_, err = s.conn.Exec(context.Background(), `
		INSERT INTO users (email, password, name) VALUES ($1, $2, $3)
	`, email, string(hashedPassword), name)
	s.Require().NoError(err, "Failed to insert test user")
}

func (s *TestSuite) loginTestUser(email, password string) {
	loginPayload := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	req, _ := http.NewRequest("POST", "/api/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", s.csrfToken)

	for _, cookie := range s.cookies {
		req.AddCookie(cookie)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Require().Equal(http.StatusNoContent, w.Code, "Expected status code 204 No Content")
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "usersession" {
			s.cookies = append(s.cookies, cookie)
		}
	}
}

func (s *TestSuite) TestCSRFToken() {
	req, _ := http.NewRequest("GET", "/api/csrf-token", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNoContent, w.Code)
	csrfToken := w.Header().Get("X-CSRF-Token")
	s.NotEmpty(csrfToken, "Expected CSRF token to be present in header")
}

func (s *TestSuite) TestRegister() {
	s.Run("should register a new user successfully", func() {
		payload := `{"email": "yozai-thinker@example.com", "password": "f205c9241173", "name": "Yozai Thinker"}`
		req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-CSRF-Token", s.csrfToken)

		for _, cookie := range s.cookies {
			req.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusNoContent, w.Code, "Expected status code 204 No Content")
	})
}

func (s *TestSuite) TestLogin() {
	s.Run("should login successfully with valid credentials", func() {
		email := "yozai-thinker@example.com"
		password := "f205c9241173"
		s.createTestUser("Yozai Thinker", email, password)

		loginPayload := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
		req, _ := http.NewRequest("POST", "/api/login", strings.NewReader(loginPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-CSRF-Token", s.csrfToken)

		for _, cookie := range s.cookies {
			req.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusNoContent, w.Code, "Expected status code 204 No Content")
	})
}

func (s *TestSuite) TestLogout() {
	s.Run("should logout successfully when user is logged in", func() {
		name := "yozai-thinker"
		email := "yozai-thinker@example.com"
		password := "f205c9241173"
		s.createTestUser(name, email, password)
		s.loginTestUser(email, password)

		req, _ := http.NewRequest("POST", "/api/logout", nil)
		req.Header.Set("X-CSRF-Token", s.csrfToken)
		for _, cookie := range s.cookies {
			req.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusNoContent, w.Code, "Expected status code 204 No Content")
	})
}

func (s *TestSuite) TestLogoutUnauthorized() {
	s.Run("should return unauthorized when user is not logged in", func() {
		req, _ := http.NewRequest("POST", "/api/logout", nil)
		req.Header.Set("X-CSRF-Token", s.csrfToken)
		for _, cookie := range s.cookies {
			req.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusUnauthorized, w.Code, "Expected status code 401 Unauthorized")
	})
}

func (s *TestSuite) TestGetUser() {
	s.Run("should return user info when logged in", func() {
		name := "yozai-thinker"
		email := "yozai-thinker@example.com"
		password := "f205c9241173"
		s.createTestUser(name, email, password)
		s.loginTestUser(email, password)

		req, _ := http.NewRequest("GET", "/api/user", nil)
		req.Header.Set("X-CSRF-Token", s.csrfToken)

		for _, cookie := range s.cookies {
			req.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)
		s.Equal(http.StatusOK, w.Code, "Expected status code 200 OK")

		expectedBody := fmt.Sprintf(`{"email":"%s","name":"%s","avatar":null,"social_accounts":[]}`, email, name)
		s.Equal(expectedBody, w.Body.String(), "Expected body to match")
	})
}

func (s *TestSuite) TestSocialAuthURL() {
	s.Run("should return valid Google auth_url", func() {
		req, _ := http.NewRequest("GET", "/api/login/social/google?redirect_uri=http://localhost/callback", nil)
		w := httptest.NewRecorder()

		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)

		var body map[string]string
		s.NoError(json.NewDecoder(w.Body).Decode(&body))

		authURL := body["auth_url"]

		expectedRegex := fmt.Sprintf(
			`^https://accounts\.google\.com/o/oauth2/auth\?access_type=offline&client_id=%s&redirect_uri=%s&response_type=code&scope=openid\+profile\+email&state=[a-f0-9]{64}$`,
			regexp.QuoteMeta(config.AppConfig.GoogleOauth2ClientID),
			regexp.QuoteMeta(url.QueryEscape("http://localhost/callback")),
		)

		match, _ := regexp.MatchString(expectedRegex, authURL)
		s.True(match, "auth_url does not match expected format:\nExpected pattern: %s\nActual: %s", expectedRegex, authURL)
	})

	s.Run("should return valid Facebook auth_url", func() {
		req, _ := http.NewRequest("GET", "/api/login/social/facebook?redirect_uri=http://localhost/callback", nil)
		w := httptest.NewRecorder()

		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)

		var body map[string]string
		s.NoError(json.NewDecoder(w.Body).Decode(&body))

		authURL := body["auth_url"]

		expectedRegex := fmt.Sprintf(
			`^https://www\.facebook\.com/v3\.2/dialog/oauth\?access_type=offline&client_id=%s&redirect_uri=%s&response_type=code&scope=email&state=[a-f0-9]{64}$`,
			regexp.QuoteMeta(config.AppConfig.FacebookOauth2ClientID),
			regexp.QuoteMeta(url.QueryEscape("http://localhost/callback")),
		)

		match, err := regexp.MatchString(expectedRegex, authURL)
		s.NoError(err)
		s.True(match, "auth_url does not match expected format:\nExpected pattern: %s\nActual: %s", expectedRegex, authURL)
	})
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
