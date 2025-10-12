// Package handlers package handlers
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
	"github.com/cheezecakee/fitrkr-athena/pkg/web"
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
	id := chi.URLParam(r, "id")

	resp, err := h.Service.GetByID(r.Context(), users.GetUserByIDReq{ID: id})
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	resp, err := h.Service.GetByUsername(r.Context(), users.GetUserByUsernameReq{Username: username})
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	resp, err := h.Service.GetByEmail(r.Context(), users.GetUserByEmailReq{Email: email})
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req users.UpdateUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.ClientError(w, http.StatusBadRequest)
		return
	}

	err := h.Service.Update(r.Context(), users.UpdateUserReq{ID: id, Username: req.Username, FirstName: req.FirstName, LastName: req.LastName, Email: req.Email})
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, "User updated")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.Service.Delete(r.Context(), users.DeleteAccountReq{ID: id})
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, "User deleted!")
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// TODO
}
