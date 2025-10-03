// Package v1
package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr-backend/internal/adapters/primary/web/v1/handlers"
	"github.com/cheezecakee/fitrkr-backend/internal/core/services/users"
)

func InitAppRoutes(userService users.UserService) http.Handler {
	r := chi.NewRouter()

	// Public routes
	r.Route("/user", func(r chi.Router) {
		r.Post("/", handlers.CreateAccount(userService))
	})

	// later you could add more groups like workouts, auth, etc.

	return r
}
