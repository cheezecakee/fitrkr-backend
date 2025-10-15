package handlers

import (
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/auth"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

type HandlerResgistry struct {
	UserHandler *UserHandler
	AuthHandler *AuthHandler
}

func NewHandlerRegistry(userService users.UserService, authService auth.AuthService) *HandlerResgistry {
	return &HandlerResgistry{
		UserHandler: NewUserHandler(userService),
		AuthHandler: NewAuthHandler(authService),
	}
}
