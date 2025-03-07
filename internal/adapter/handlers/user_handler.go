package handlers

import (
	"net/http"

	"github.com/Joe5451/go-oauth2-server/internal/application/ports/in"
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

func (h *UserHandler) Register(c *gin.Context) {
	json := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Username string `json:"name" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "請求參數錯誤，請檢查輸入",
			"error":   err.Error(),
		})
		return
	}

	err := h.usecase.Register(in.RegisterUserRequest{
		Email:    json.Email,
		Password: json.Password,
		Username: json.Username,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "發生錯誤",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
