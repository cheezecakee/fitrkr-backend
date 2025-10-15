package middleware

import (
	"net/http"
	"slices"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func RequireAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logr.Get().Infof("RequireAdmin: Context keys available: %+v", r.Context())

			currentUser, ok := r.Context().Value(UserKey).(*user.User)
			logr.Get().Infof("RequireAdmin: userCtx extraction - ok: %v, userCtx: %+v", ok, currentUser)

			if !ok || currentUser == nil {
				logr.Get().Infof("RequireAdmin: Failed to get user from context")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			logr.Get().Infof("RequireAdmin: User found: %s, Roles: %v, Role type: %T", currentUser.Username, currentUser.Roles, currentUser.Roles)

			if len(currentUser.Roles) == 0 {
				logr.Get().Infof("RequireAdmin: User has no roles assigned")
			}

			// Fix: Compare Role types, not strings
			if slices.Contains(currentUser.Roles, user.RoleAdmin) {
				logr.Get().Infof("RequireAdmin: Admin access granted")
				next.ServeHTTP(w, r)
				return
			}

			logr.Get().Infof("RequireAdmin: User is not admin. Expected 'admin', got roles: %v", currentUser.Roles)
			http.Error(w, "forbidden: admin only", http.StatusForbidden)
		})
	}
}

func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentUser, ok := r.Context().Value(UserKey).(*user.User)
			if !ok || currentUser == nil {
				logr.Get().Infof("RequireAnyRole: Failed to get user from context")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			logr.Get().Infof("RequireAnyRole: Checking if user %s has any of roles: %v (user has: %v)",
				currentUser.Username, roles, currentUser.Roles)

			// Fix: Convert string roles to Role types for comparison
			for _, requiredRoleStr := range roles {
				requiredRole := user.Role(requiredRoleStr)
				if slices.Contains(currentUser.Roles, requiredRole) {
					logr.Get().Infof("RequireAnyRole: Access granted - user has role: %s", requiredRole)
					next.ServeHTTP(w, r)
					return
				}
			}

			logr.Get().Infof("RequireAnyRole: Access denied - user lacks required roles: %v", roles)
			http.Error(w, "forbidden: insufficient privileges", http.StatusForbidden)
		})
	}
}
