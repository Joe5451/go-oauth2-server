package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Joe5451/go-oauth2-server/internal/application/ports/in"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	ErrValidation   = errors.New("Validation error")
	ErrUnauthorized = errors.New("Requires authentication.")
)

type UserHandler struct {
	usecase in.UserUsecase
}

func NewUserHandler(usecase in.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) CSRFToken(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Register(c *gin.Context) {
	json := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(fmt.Errorf("%w: %v", ErrValidation, err.Error()))
		return
	}

	err := h.usecase.Register(in.RegisterUserRequest{
		Email:    json.Email,
		Password: json.Password,
		Name:     json.Name,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) LoginWithEmail(c *gin.Context) {
	json := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(fmt.Errorf("%w: %v", ErrValidation, err.Error()))
		return
	}

	user, err := h.usecase.AuthenticateUser(json.Email, json.Password)
	if err != nil {
		c.Error(err)
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("user_id")

	if v == nil {
		c.Error(ErrUnauthorized)
		return
	}

	userID := v.(int64)
	user, err := h.usecase.GetUser(userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":  user.Email,
		"name":   user.Name,
		"avatar": user.Avatar,
	})
}
