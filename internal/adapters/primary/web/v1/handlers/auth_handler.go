// Package handlers package handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/adapters/primary/web/middleware"
	"github.com/cheezecakee/fitrkr-athena/internal/adapters/secondary/external/jwt"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/auth"
	"github.com/cheezecakee/fitrkr-athena/pkg/web"
)

type AuthHandler struct {
	Service    auth.AuthService
	jwtManager jwt.JWT
}

func NewAuthHandler(service auth.AuthService, jwtManager jwt.JWT) *AuthHandler {
	return &AuthHandler{Service: service, jwtManager: jwtManager}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.ClientError(w, http.StatusBadRequest)
		return
	}

	resp, err := h.Service.Login(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}

	token, err := h.jwtManager.MakeJWT(resp.UserID, resp.Roles)
	if err != nil {
		web.ServerError(w, err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: false, // Set to true in production with HTTPS
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(15 * time.Minute), // Adjust as needed
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken, // access token
		Path:     "/auth/refresh",
		HttpOnly: false, // Set to true in production with HTTPS
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // Adjust as needed
	})

	w.Header().Set("Content-Type", "application/json")

	web.Response(w, http.StatusOK, "Login successful")
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := middleware.ExtractToken(r, middleware.RefreshToken)
	if err != nil {
		web.ServerError(w, err)
	}

	refresh := auth.RefreshReq{Token: token}

	resp, err := h.Service.Refresh(r.Context(), refresh)
	if err != nil {
		web.ServerError(w, err)
	}

	accessToken, err := h.jwtManager.MakeJWT(resp.UserID, resp.Roles)
	if err != nil {
		web.ServerError(w, err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: false, // Set to true in production with HTTPS
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(15 * time.Minute), // Adjust as needed
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.Token, // access token
		Path:     "/auth/refresh",
		HttpOnly: false, // Set to true in production with HTTPS
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // Adjust as needed
	})

	w.Header().Set("Content-Type", "application/json")

	web.Response(w, http.StatusOK, "Login successful")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req auth.RevokeTokenReq

	token, err := middleware.ExtractToken(r, middleware.RefreshToken)
	if err != nil {
		web.ServerError(w, err)
	}

	req.Token = token

	err = h.Service.Revoke(r.Context(), req)
	if err != nil {
		web.ServerError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Deletes the cookie
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh",
		MaxAge:   -1,
		HttpOnly: true,
	})

	logr.Get().Info("User logged out successfully!")

	w.WriteHeader(http.StatusNoContent)
}
