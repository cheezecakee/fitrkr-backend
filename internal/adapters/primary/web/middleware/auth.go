// Pacakge middleware
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cheezecakee/logr"
)

func extractToken(r *http.Request) (string, error) {
	var token string

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
		logr.Get().Infof("Found Bearer token in Authorization header")
		return token, nil
	}

	cookieNames := []string{"session"}

	for _, cookieName := range cookieNames {
		if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
			token = cookie.Value
			logr.Get().Infof("Found JWT token in cookie: %s", cookieName)
			return token, nil
		}
	}

	logr.Get().Infof("No JWT token found in Authorization header or cookies")
	return "", errors.ErrTokenNotFound
}

func IsAuthenticated() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log request details for debugging remove in production
			logr.Get().Infof("Auth check for %s %s", r.Method, r.URL.Path)
			logr.Get().Infof("Request headers: %v", r.Header)

			// Extract token from either Bearer header or cookies
			token, err := extractToken(r)
			if err != nil {
				logr.Get().Infof("Token extraction failed: %v", err)
				http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			logr.Get().Infof("Extracted token (first 20 chars): %s...",
				func() string {
					if len(token) > 20 {
						return token[:20]
					}
					return token
				}())

			userID, err := h.app.JWTManager.ValidateJWT(token)
			if err != nil {
				logr.Get().Infof("JWT validation failed: %v", err)
				http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			logr.Get().Infof("JWT validation successful for user ID: %s", userID)

			user, err := h.app.UserSvc.GetUserByID(r.Context(), userID)
			if err != nil {
				logr.Get().Infof("Failed to get user by ID %s: %v", userID, err)
				http.Error(w, "unauthorized: user not found", http.StatusUnauthorized)
				return
			}

			logr.Get().Infof("User found: %s with roles: %v", user.Username, user.Roles)

			ctx := context.WithValue(r.Context(), UserKey, &user)
			ctx = context.WithValue(ctx, UserIDKey, user.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
