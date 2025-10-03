// Package web
package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	v1 "github.com/cheezecakee/fitrkr-backend/internal/adapters/primary/web/v1"
	"github.com/cheezecakee/fitrkr-backend/internal/core/services/users"
)

type App struct {
	chi         *chi.Mux
	userService users.UserService
	port        int
}

func NewApp(userService users.UserService, opts ...AppOption) *App {
	app := &App{
		userService: userService,
		port:        8000,
		chi:         chi.NewRouter(),
	}

	for _, applyOption := range opts {
		applyOption(app)
	}

	app.chi.Mount("/api/v1", v1.InitAppRoutes(userService))

	return app
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.port)
	return http.ListenAndServe(addr, a.chi)
}
