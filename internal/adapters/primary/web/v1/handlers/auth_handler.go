// Package handlers package handlers
package handlers

import (
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/auth"
)

type AuthHandler struct {
	Service auth.AuthService
}

func NewAuthHandler(service auth.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}
