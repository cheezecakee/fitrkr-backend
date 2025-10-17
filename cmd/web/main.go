// Package main
package main

import (
	"os"

	"github.com/cheezecakee/logr"
	"github.com/joho/godotenv"

	"github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/secondary/db/postgres"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/secondary/external/jwt"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/auth"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

func main() {
	logr.Init(&logr.PlainTextFormatter{}, logr.LevelInfo, nil)

	if err := godotenv.Load(); err != nil {
		logr.Get().Errorf("No .env file found: %v", err)
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		logr.Get().Error("JWT_SECRET environment variable not set")
	}
	jwtManager := jwt.NewJWTManager(secretKey, 0)

	db := postgres.NewPostgresConn()
	defer db.Close()

	userRepo, err := postgres.NewUserRepo(db)
	if err != nil {
		logr.Get().Errorf("failed to init postgres user repo: %v", err)
	}
	authRepo, err := postgres.NewAuthRepo(db)
	if err != nil {
		logr.Get().Errorf("failed to init postgres auth repo: %v", err)
	}

	userService := users.NewService(userRepo)
	authService := auth.NewService(authRepo, userRepo)

	server := web.NewApp(
		userService,
		authService,
		jwtManager,
		web.WithPort(8000))

	logr.Get().Info("Starting server...")
	server.Run()
}
