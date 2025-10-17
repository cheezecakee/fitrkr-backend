// Package web
package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web/middleware"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web/v1"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web/v1/handlers"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/secondary/external/jwt"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/auth"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

type App struct {
	chi        *chi.Mux
	registry   *handlers.HandlerResgistry
	middleware *middleware.Middleware
	jwtManager jwt.JWT
	port       int
}

func NewApp(userService users.UserService, authService auth.AuthService, jwtManager jwt.JWT, opts ...AppOption) *App {
	app := &App{
		port:       8000,
		chi:        chi.NewRouter(),
		jwtManager: jwtManager,
	}

	for _, applyOption := range opts {
		applyOption(app)
	}

	app.middleware = middleware.NewMiddleware(jwtManager)
	app.registry = handlers.NewHandlerRegistry(userService, authService, jwtManager, *app.middleware)

	fs := http.FileServer(http.Dir("internal/adapters/primary/web/docs"))

	app.chi.Use(app.registry.CORS)
	app.chi.Handle("/api/v1/docs/*", http.StripPrefix("/api/v1/docs/", fs))

	app.chi.Mount("/api/v1", v1.RegisterRoutes(app.registry))

	return app
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.port)
	return http.ListenAndServe(addr, a.chi)
}
