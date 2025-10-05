// Package handlers package handlers
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cheezecakee/fitrkr-backend/internal/core/services/users"
	"github.com/cheezecakee/fitrkr-backend/pkg/web"
)

type UserHandler struct {
	Service users.UserService
}

func NewUserHandler(service users.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req users.CreateAccountReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.ClientError(w, http.StatusBadRequest)
		return
	}

	resp, err := h.Service.CreateAccount(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusCreated, resp)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// TODO
}
