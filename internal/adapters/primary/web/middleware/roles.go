package middleware

import (
	"net/http"

	"github.com/cheezecakee/logr"

	webctx "github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web/context"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/secondary/external/jwt"
	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func RequireAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUser, ok := r.Context().Value(webctx.AuthenticatedUserKey).(*jwt.AuthenticatedUser)
			logr.Get().Debugf("RequireAdmin: userCtx extraction - ok: %v, userCtx: %+v", ok, authUser)
			if !ok || authUser == nil {
				logr.Get().Errorf("failed to authorize, admin only")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if authUser.Roles.Contains(user.RoleAdmin) {
				logr.Get().Info("admin access granted")
				next.ServeHTTP(w, r)
				return
			}

			logr.Get().Error("user access denied, admin only")
			http.Error(w, "forbidden: admin only", http.StatusForbidden)
		})
	}
}

func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUser, ok := r.Context().Value(webctx.AuthenticatedUserKey).(*jwt.AuthenticatedUser)
			if !ok || authUser == nil {
				logr.Get().Error("failed to authorize, login required")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if authUser.Roles.ContainsAny(user.RoleAdmin, user.RoleModerator, user.RoleUser) {
				logr.Get().Info("access granted")
				next.ServeHTTP(w, r)
				return
			}

			logr.Get().Error("user acces denied, login required")
			http.Error(w, "forbidden: insufficient privileges", http.StatusForbidden)
		})
	}
}
