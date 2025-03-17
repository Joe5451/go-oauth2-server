package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
}

func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{}
}

func (h *TemplateHandler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title":   "Login",
		"showNav": false,
	})
}

func (h *TemplateHandler) SocialLinks(c *gin.Context) {
	c.HTML(http.StatusOK, "social_links.tmpl", gin.H{
		"title":   "User Social Links",
		"showNav": true,
	})
}
