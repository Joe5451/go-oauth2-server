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

func NewRouter(
	userHandler *handlers.UserHandler,
	templateHandler *handlers.TemplateHandler,
) *gin.Engine {
	router := gin.Default()

	// Set up a redis session
	store, _ := redis.NewStore(
		10,
		"tcp",
		fmt.Sprintf("%s:%s", config.AppConfig.RedisHost, config.AppConfig.RedisPort),
		config.AppConfig.RedisPassword,
		[]byte(config.AppConfig.RedisSecret),
	)

	// API
	router.Static("/uploads", "./uploads")
	{
		api := router.Group("/api")

		// Set up session
		api.Use(sessions.Sessions("usersession", store))

		// Set up error handler
		api.Use(middlewares.InitErrorHandler())

		// setup csrf middleware
		api.Use(middlewares.CSRF())
		api.Use(middlewares.CSRFToken())

		api.GET("/csrf-token", userHandler.CSRFToken)
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.LoginWithEmail)
		api.POST("/logout", userHandler.Logout)
		api.GET("/user", userHandler.GetUser)
		api.PATCH("/user/avatar", userHandler.UpdateUserAvatar)

		api.GET("/login/social/:provider", userHandler.SocialAuthURL)
		api.POST("/login/social/callback", userHandler.SocialAuthCallback)

		api.GET("/auth/social/:provider/link/url", userHandler.SocialAuthUrlForLinkingExistingUser)
		api.POST("/auth/social/link", userHandler.LinkUserWithSocialAccount)

		api.POST("/user/link/:provider", userHandler.LinkSocialAccount)
		api.DELETE("/user/unlink/:provider", userHandler.UnlinkSocialAccount)
	}

	// Template
	router.Static("/assets", "./web/assets")
	router.LoadHTMLGlob("web/templates/*.tmpl")
	{
		template := router.Group("/template")
		template.GET("/login", templateHandler.Login)
		template.GET("/user/social-links", templateHandler.SocialLinks)
	}

	return router
}
