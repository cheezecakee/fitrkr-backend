// Package webctx
package webctx

type ContextKey string

const (
	UserIDKey            ContextKey = "user_id"
	AuthenticatedUserKey ContextKey = "authenticated_user"
)
