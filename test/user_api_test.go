package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
