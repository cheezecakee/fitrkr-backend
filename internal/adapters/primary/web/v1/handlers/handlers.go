package handlers

import (
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web/middleware"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/secondary/external/jwt"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/auth"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

type HandlerResgistry struct {
	UserHandler *UserHandler
	AuthHandler *AuthHandler
	JwtManager  jwt.JWT
	*middleware.Middleware
}

func NewHandlerRegistry(userService users.UserService, authService auth.AuthService, jwtManager jwt.JWT, middleware middleware.Middleware) *HandlerResgistry {
	return &HandlerResgistry{
		UserHandler: NewUserHandler(userService),
		AuthHandler: NewAuthHandler(authService, jwtManager),
		JwtManager:  jwtManager,
		Middleware:  &middleware,
	}
}
