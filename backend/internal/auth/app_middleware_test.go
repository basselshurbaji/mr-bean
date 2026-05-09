package auth_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/auth"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

func TestAppMiddleware_MissingHeader(t *testing.T) {
	mw := auth.AppMiddleware(&mockAppTokenService{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestAppMiddleware_NoBearerPrefix(t *testing.T) {
	mw := auth.AppMiddleware(&mockAppTokenService{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token abc123")
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAppMiddleware_ValidateError(t *testing.T) {
	svc := &mockAppTokenService{valErr: errors.New("token revoked")}
	mw := auth.AppMiddleware(svc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer some.raw.token")
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestAppMiddleware_ValidToken_SetsUserID(t *testing.T) {
	svc := &mockAppTokenService{userID: "user-123"}
	mw := auth.AppMiddleware(svc)

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
	req.Header.Set("Authorization", "Bearer valid.app.token")
	rec := httptest.NewRecorder()
	mw(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if capturedID != "user-123" {
		t.Errorf("expected user-123, got %s", capturedID)
	}
}

func TestAppMiddleware_RealJWT_Integration(t *testing.T) {
	tokenSvc := auth.NewTokenService("test-secret", time.Minute, time.Hour)
	raw, err := tokenSvc.GenerateAppToken("user-456", "db-record-id")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	// Service mock that validates the JWT via real token service and returns the user ID.
	svc := &mockAppTokenService{userID: "user-456"}
	mw := auth.AppMiddleware(svc)

	var capturedID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := principal.UserIDFromContext(r.Context())
		capturedID = id
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+raw)
	rec := httptest.NewRecorder()
	mw(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if capturedID != "user-456" {
		t.Errorf("expected user-456, got %s", capturedID)
	}
}
