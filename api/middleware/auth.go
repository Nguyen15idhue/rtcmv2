package middleware

import (
	"net/http"
	"os"
	"strings"
)

type AuthMiddleware struct {
	apiKey string
}

func NewAuthMiddleware() *AuthMiddleware {
	key := os.Getenv("API_KEY")
	if key == "" {
		key = "default-api-key-change-me"
	}
	return &AuthMiddleware{apiKey: key}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skipAuth(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			http.Error(w, `{"error":"missing_api_key","message":"X-API-Key header is required"}`, http.StatusUnauthorized)
			return
		}

		if apiKey != m.apiKey {
			http.Error(w, `{"error":"invalid_api_key","message":"Invalid API key"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) SetAPIKey(key string) {
	m.apiKey = key
}

func (m *AuthMiddleware) GetAPIKey() string {
	return m.apiKey
}

func skipAuth(path string) bool {
	if path == "/" {
		return true
	}
	if path == "/api/health" {
		return true
	}
	if path == "/api/stream" {
		return true
	}
	if path == "/api/system" {
		return true
	}
	if path == "/api/stations" {
		return true
	}
	if path == "/api/station" {
		return true
	}
	if strings.HasPrefix(path, "/api/station/") {
		return true
	}
	if path == "/api/casters" {
		return true
	}
	if path == "/api/caster" {
		return true
	}
	if strings.HasPrefix(path, "/api/caster/") {
		return true
	}

	return false
}
