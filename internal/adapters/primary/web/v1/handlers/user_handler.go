// Package handlers package handlers
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cheezecakee/fitrkr-backend/internal/core/services/users"
	"github.com/cheezecakee/fitrkr-backend/pkg/web"
)

func CreateAccount(userService users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req users.CreateAccountReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			web.ClientError(w, http.StatusBadRequest)
			return
		}

		resp, err := userService.CreateAccount(r.Context(), req)
		if err != nil {
			web.ServerError(w, err)
			return
		}

		web.Response(w, http.StatusCreated, resp)
	}
}
