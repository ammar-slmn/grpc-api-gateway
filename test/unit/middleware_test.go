package loadbalancer_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"grpc-api-gateway/pkg/middleware"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		config         middleware.AuthConfig
		setupRequest   func(*http.Request)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid API Key",
			config: middleware.AuthConfig{
				Enabled: true,
				APIKeys: map[string]string{"valid-key": "test-user"},
			},
			setupRequest: func(r *http.Request) {
				r.Header.Set("X-API-Key", "valid-key")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing API Key",
			config: middleware.AuthConfig{
				Enabled: true,
				APIKeys: map[string]string{"valid-key": "test-user"},
			},
			setupRequest:   func(r *http.Request) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing API key\n",
		},
		{
			name: "Invalid API Key",
			config: middleware.AuthConfig{
				Enabled: true,
				APIKeys: map[string]string{"valid-key": "test-user"},
			},
			setupRequest: func(r *http.Request) {
				r.Header.Set("X-API-Key", "invalid-key")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid API key\n",
		},
		{
			name: "Auth Disabled",
			config: middleware.AuthConfig{
				Enabled: false,
				APIKeys: map[string]string{},
			},
			setupRequest:   func(r *http.Request) {},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			middleware := middleware.NewAuthMiddleware(tt.config)(nextHandler)

			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupRequest(req)
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" {
				if body := w.Body.String(); body != tt.expectedBody {
					t.Errorf("Expected body %q, got %q", tt.expectedBody, body)
				}
			}
		})
	}
}
