package user

import (
	"errors"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
)

var ErrInvalidRole = errors.New("invalid role")

func NewRole(roles []string) ([]Role, error) {
	var Roles []Role

	// If no role is specified on account creation defaults to user
	if roles == nil {
		Roles = append(Roles, RoleUser)
		return Roles, nil
	}

	for _, role := range roles {
		switch role {
		case "user":
			Roles = append(Roles, RoleUser)
		case "admin":
			Roles = append(Roles, RoleAdmin)
		case "moderator":
			Roles = append(Roles, RoleModerator)
		default:
			return nil, ErrInvalidRole
		}
	}

	return Roles, nil
}

// Helper functions

func RolesToStrings(roles []Role) []string {
	strings := make([]string, len(roles))
	for i, role := range roles {
		strings[i] = string(role)
	}
	return strings
}

func StringsToRoles(strings []string) []Role {
	roles := make([]Role, len(strings))
	for i, s := range strings {
		roles[i] = Role(s)
	}
	return roles
}
