package middleware

import (
	"net/http"
	"slices"

	"github.com/cheezecakee/logr"
)

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow multiple origins - add your actual Flutter app URLs
		allowedOrigins := []string{
			"http://localhost:8000", // OpenApi port
			"http://localhost",
			"capacitor://localhost", // For Capacitor apps
			"ionic://localhost",     // For Ionic apps
		}

		origin := r.Header.Get("Origin")

		// Debug logging - remove in production
		logr.Get().Infof("Request Origin: %s", origin)
		logr.Get().Infof("Request Method: %s", r.Method)
		logr.Get().Infof("Request Headers: %v", r.Header)

		if slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			logr.Get().Infof("Origin allowed: %s", origin)
		} else if origin == "" {
			// For local API requests (same origin, like Stoplight UI)
			// Allow localhost explicitly instead of "*"
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
			logr.Get().Info("No Origin header, using default localhost:8000")
		} else {
			logr.Get().Infof("Origin not allowed: %s", origin)
			// Might want to still set some CORS headers for the error response
			w.Header().Set("Access-Control-Allow-Origin", "null")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours

		if r.Method == http.MethodOptions {
			logr.Get().Infof("Handling OPTIONS preflight request from: %s", origin)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
