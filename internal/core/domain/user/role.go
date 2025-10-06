package user

import (
	"errors"
)

type (
	Role  string
	Roles []Role
)

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
)

var ErrInvalidRole = errors.New("invalid role")

func NewRoles(strs []string) (Roles, error) {
	if strs == nil {
		return Roles{RoleUser}, nil
	}
	if len(strs) == 0 {
		return Roles{RoleUser}, nil
	}

	roles := make(Roles, 0, len(strs))
	for _, s := range strs {
		switch s {
		case string(RoleUser):
			roles = append(roles, RoleUser)
		case string(RoleAdmin):
			roles = append(roles, RoleAdmin)
		case string(RoleModerator):
			roles = append(roles, RoleModerator)
		default:
			return nil, ErrInvalidRole
		}
	}

	return roles, nil
}

// Helper functions

func (r Roles) ToStrings() []string {
	out := make([]string, len(r))
	for i, role := range r {
		out[i] = string(role)
	}
	return out
}

func StringsToRoles(strs []string) Roles {
	out := make(Roles, len(strs))
	for i, s := range strs {
		out[i] = Role(s)
	}
	return out
}
