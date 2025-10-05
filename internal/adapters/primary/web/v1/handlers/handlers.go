package handlers

import (
	"github.com/cheezecakee/fitrkr-backend/internal/core/services/users"
)

type HandlerResgistry struct {
	UserHandler *UserHandler
}

func NewHandlerRegistry(userService users.UserService) *HandlerResgistry {
	return &HandlerResgistry{
		UserHandler: NewUserHandler(userService),
	}
}
