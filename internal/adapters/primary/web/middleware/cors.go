package middleware

import (
	"net/http"
	"slices"

	"github.com/cheezecakee/logr"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow multiple origins - add your actual Flutter app URLs
		allowedOrigins := []string{
			"http://localhost:5173", // Your Svelte app
			"http://localhost:3000", // Common Flutter web port
			"http://localhost:8080", // Another common Flutter web port
			"http://10.0.2.2:3000",  // Android emulator
			"http://127.0.0.1:3000", // Alternative localhost
			"https://receiver-consistently-exchange-women.trycloudflare.com", // Your Cloudflare tunnel
			"http://localhost",      // For mobile apps
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
			// For mobile apps that don't send Origin header, allow them
			w.Header().Set("Access-Control-Allow-Origin", "*")
			logr.Get().Infof("No origin header, allowing all")
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
