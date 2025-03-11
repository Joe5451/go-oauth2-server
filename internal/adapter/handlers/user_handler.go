package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Joe5451/go-oauth2-server/internal/application/ports/in"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase in.UserUsecase
}

func NewUserHandler(usecase in.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func handleError(c *gin.Context, err error) {
	if errInfo, exists := errorMap[err]; exists {
		c.JSON(errInfo.httpCode, gin.H{
			"code":    errInfo.errorCode,
			"message": errInfo.message,
		})
	} else if errors.Is(err, ErrValidation) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "VALIDATION_ERROR",
			"message": err.Error(),
		})
	} else {
		log.Printf("Register failed. error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred.",
		})
	}
	return
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
		handleError(c, fmt.Errorf("%w: %v", ErrValidation, err.Error()))
		return
	}

	err := h.usecase.Register(in.RegisterUserRequest{
		Email:    json.Email,
		Password: json.Password,
		Name:     json.Name,
	})

	if err != nil {
		handleError(c, err)
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
		handleError(c, fmt.Errorf("%w: %v", ErrValidation, err.Error()))
		return
	}

	user, err := h.usecase.AuthenticateUser(json.Email, json.Password)

	if err != nil {
		handleError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.Status(http.StatusNoContent)
}
