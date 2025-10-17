// Package v1
package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web/v1/handlers"
)

func RegisterRoutes(resgitry *handlers.HandlerResgistry) http.Handler {
	r := chi.NewRouter()

	routes := map[string]http.Handler{
		"/user": SetupUserRoutes(resgitry),
		"/auth": SetupAuthRoutes(resgitry),
	}

	// Mount the versioned routes
	for path, handler := range routes {
		r.Mount(path, handler)
	}

	return r
}

func SetupUserRoutes(registry *handlers.HandlerResgistry) http.Handler {
	r := chi.NewRouter()

	r.Post("/", registry.UserHandler.CreateAccount)

	r.Group(func(r chi.Router) {
		r.Use(registry.IsAuthenticated())
		r.Get("/username/{username}", registry.UserHandler.GetUserByUsername) // Make this public (view profile)
		r.Get("/", registry.UserHandler.GetUserByID)
		r.Get("/email/{email}", registry.UserHandler.GetUserByEmail)
		r.Put("/", registry.UserHandler.UpdateUser)
		r.Delete("/", registry.UserHandler.DeleteUser)

		r.Get("/subscription", registry.UserHandler.GetSubscription)
		r.Get("/settings", registry.UserHandler.GetSettings)
		r.Get("/stats", registry.UserHandler.GetStats)

		r.Put("/settings", registry.UserHandler.UpdateSettings)
		r.Put("/stats/body", registry.UserHandler.UpdateBodyMetrics)
		r.Put("/subscription/plan", registry.UserHandler.UpgradePlan)
		r.Put("/subscription/payment", registry.UserHandler.RecordPayment)
		r.Put("/subscription/cancel", registry.UserHandler.CancelSubscription)
		r.Put("/subscription/trial", registry.UserHandler.StartTrial)
	})
	return r
}

func SetupAuthRoutes(registry *handlers.HandlerResgistry) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", registry.AuthHandler.Login)
	r.Post("/refresh", registry.AuthHandler.Refresh)
	r.Group(func(r chi.Router) {
		r.Use(registry.IsAuthenticated())
		r.Post("/logout", registry.AuthHandler.Logout)
	})
	return r
}
