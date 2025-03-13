package http

import (
	"fmt"

	"github.com/Joe5451/go-oauth2-server/internal/adapter/handlers"
	"github.com/Joe5451/go-oauth2-server/internal/config"
	"github.com/Joe5451/go-oauth2-server/internal/http/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	// Set up a redis session
	store, _ := redis.NewStore(
		10,
		"tcp",
		fmt.Sprintf("%s:%s", config.AppConfig.RedisHost, config.AppConfig.RedisPort),
		config.AppConfig.RedisPassword,
		[]byte(config.AppConfig.RedisSecret),
	)

	router.Use(sessions.Sessions("usersession", store))

	// Set up error handler
	router.Use(middlewares.InitErrorHandler())

	// setup csrf middleware
	router.Use(middlewares.CSRF())
	router.Use(middlewares.CSRFToken())

	router.GET("/csrf-token", userHandler.CSRFToken)
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.LoginWithEmail)
	router.POST("/logout", userHandler.Logout)
	router.GET("/user", userHandler.GetUser)

	router.GET("/login/social/:provider", userHandler.SocialAuthURL)
	router.POST("/login/social/callback", userHandler.SocialAuthCallback)

	router.GET("/auth/social/:provider/link/url", userHandler.LinkSocialAuthUrl)
	router.POST("/auth/social/link", userHandler.LinkUserWithSocialAccount)

	// router.GET("/auth/:provider/url", userHandler.GenerateAuthUrl)
	// router.POST("/auth/:provider/callback", userHandler.HandleOAuth2Callback)
	// router.PATCH("/user", userHandler.UpdateUser)

	return router
}
