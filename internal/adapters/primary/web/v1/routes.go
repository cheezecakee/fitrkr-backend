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
	r.Get("/username/{username}", registry.UserHandler.GetUserByUsername)
	r.Get("/{id}", registry.UserHandler.GetUserByID)
	r.Get("/email/{email}", registry.UserHandler.GetUserByEmail)
	r.Put("/{id}", registry.UserHandler.UpdateUser)
	r.Delete("/{id}", registry.UserHandler.DeleteUser)

	r.Get("/{id}/subscription", registry.UserHandler.GetSubscription)
	r.Get("/{id}/settings", registry.UserHandler.GetSettings)
	r.Get("/{id}/stats", registry.UserHandler.GetStats)

	r.Put("/{id}/settings", registry.UserHandler.UpdateSettings)
	r.Put("/{id}/stats/body", registry.UserHandler.UpdateBodyMetrics)
	r.Put("/{id}/subscription/plan", registry.UserHandler.UpgradePlan)
	r.Put("/{id}/subscription/payment", registry.UserHandler.RecordPayment)
	r.Put("/{id}/subscription/cancel", registry.UserHandler.CancelSubscription)
	r.Put("/{id}/subscription/trial", registry.UserHandler.StartTrial)
	return r
}
