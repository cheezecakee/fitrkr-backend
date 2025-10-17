// Package middleware
package middleware

import (
	"context"
	"net/http"

	"github.com/cheezecakee/logr"

	webctx "github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web/context"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/secondary/external/jwt"
)

type Middleware struct {
	jwtManager jwt.JWT
}

func NewMiddleware(jwtManager jwt.JWT) *Middleware {
	return &Middleware{
		jwtManager: jwtManager,
	}
}

func (m *Middleware) IsAuthenticated() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log request details for debugging remove in production
			logr.Get().Debugf("auth check for %s %s", r.Method, r.URL.Path)
			logr.Get().Debugf("request headers: %v", r.Header)

			// Extract token from either Bearer header or cookies
			token, err := ExtractToken(r, Session)
			if err != nil {
				logr.Get().Errorf("token extraction failed: %v", err)
				http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			authUser, err := m.jwtManager.ValidateJWT(token)
			if err != nil {
				logr.Get().Infof("JWT validation failed: %v", err)
				http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			logr.Get().Debugf("JWT validation successful for user ID: %s", authUser.UserID)

			ctx := context.WithValue(r.Context(), webctx.AuthenticatedUserKey, authUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
