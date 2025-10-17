// Package handlers package handlers
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	webctx "github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web/context"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/secondary/external/jwt"
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
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}

	resp, err := h.Service.GetByID(r.Context(), users.GetUserByIDReq{ID: user.UserID.String()})
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

// GetUserByEmail Make this admin only later
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
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}

	var req users.UpdateUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.ClientError(w, http.StatusBadRequest)
		return
	}

	req.ID = user.UserID.String()

	err = h.Service.Update(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, "User updated")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}

	err = h.Service.Delete(r.Context(), users.DeleteAccountReq{ID: user.UserID.String()})
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, "User deleted!")
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *UserHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}

	resp, err := h.Service.GetSubscription(r.Context(), users.GetSubscriptionReq{ID: user.UserID.String()})
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, resp.Subscription)
}

func (h *UserHandler) UpgradePlan(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}

	var req users.UpgradePlanReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.ClientError(w, http.StatusBadRequest)
		return
	}

	req.UserID = user.UserID.String()
	err = h.Service.UpgradePlan(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}
	web.Response(w, http.StatusOK, "User plan upgraded")
}

func (h *UserHandler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}
	var req users.RecordPaymentReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.ClientError(w, http.StatusBadRequest)
		return
	}

	req.UserID = user.UserID.String()
	err = h.Service.RecordPayment(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}
	web.Response(w, http.StatusOK, "User payment recorded")
}

func (h *UserHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}
	var req users.CancelSubscriptionReq

	req.UserID = user.UserID.String()
	err = h.Service.CancelSubscription(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}
	web.Response(w, http.StatusOK, "User subscription cancelled")
}

func (h *UserHandler) StartTrial(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}
	var req users.StartTrialReq

	req.UserID = user.UserID.String()
	err = h.Service.StartTrial(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}
	web.Response(w, http.StatusOK, "User trail started")
}

func (h *UserHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}

	resp, err := h.Service.GetSettings(r.Context(), users.GetSettingsReq{ID: user.UserID.String()})
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, resp.Settings)
}

func (h *UserHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}

	resp, err := h.Service.GetStats(r.Context(), users.GetStatsReq{ID: user.UserID.String()})
	if err != nil {
		web.ServerError(w, err)
		return
	}

	web.Response(w, http.StatusOK, resp.Stats)
}

func (h *UserHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}

	var req users.UpdateSettingsReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.ClientError(w, http.StatusBadRequest)
		return
	}

	req.UserID = user.UserID.String()
	err = h.Service.UpdateSettings(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}
	web.Response(w, http.StatusOK, "User settings updated")
}

func (h *UserHandler) UpdateBodyMetrics(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		web.ClientError(w, http.StatusUnauthorized)
		return
	}

	var req users.UpdateBodyMetricsReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.ClientError(w, http.StatusBadRequest)
		return
	}

	req.UserID = user.UserID.String()
	err = h.Service.UpdateBodyMetrics(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}
	web.Response(w, http.StatusOK, "User body metrics updated")
}

func getUser(ctx context.Context) (*jwt.AuthenticatedUser, error) {
	user := ctx.Value(webctx.AuthenticatedUserKey).(*jwt.AuthenticatedUser)
	if user.UserID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized user")
	}

	return user, nil
}
