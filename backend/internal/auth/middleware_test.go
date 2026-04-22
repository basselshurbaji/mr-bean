package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/auth"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

func TestMiddleware_MissingHeader(t *testing.T) {
	svc := auth.NewTokenService("secret", time.Minute, time.Hour)
	mw := auth.Middleware(svc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestMiddleware_NoBearerPrefix(t *testing.T) {
	svc := auth.NewTokenService("secret", time.Minute, time.Hour)
	mw := auth.Middleware(svc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token abc123")
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestMiddleware_InvalidToken(t *testing.T) {
	svc := auth.NewTokenService("secret", time.Minute, time.Hour)
	mw := auth.Middleware(svc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestMiddleware_ValidToken_SetsUserID(t *testing.T) {
	svc := auth.NewTokenService("secret", time.Minute, time.Hour)
	mw := auth.Middleware(svc)

	token, err := svc.GenerateAccessToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	var capturedID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := principal.UserIDFromContext(r.Context())
		if !ok {
			t.Error("expected user ID in context")
		}
		capturedID = id
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mw(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if capturedID != "user-123" {
		t.Errorf("expected user-123, got %s", capturedID)
	}
}

func TestMiddleware_RefreshTokenRejected(t *testing.T) {
	svc := auth.NewTokenService("secret", time.Minute, time.Hour)
	mw := auth.Middleware(svc)

	refresh, err := svc.GenerateRefreshToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+refresh)
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
