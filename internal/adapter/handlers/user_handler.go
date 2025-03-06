package handlers

import (
	"github.com/Joe5451/go-oauth2-server/internal/application/ports/in"
)

type UserHandler struct {
	usecase in.UserUsecase
}

func NewUserHandler(usecase in.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}
