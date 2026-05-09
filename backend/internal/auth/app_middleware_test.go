package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/auth"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type mockAppMiddlewareSvc struct {
	userID string
	valErr error
}

func (m *mockAppMiddlewareSvc) Create(_ context.Context, _, _ string) (string, *auth.AppToken, error) {
	return "", nil, nil
}
func (m *mockAppMiddlewareSvc) Revoke(_ context.Context, _, _ string) error { return nil }
func (m *mockAppMiddlewareSvc) Validate(_ context.Context, _ string) (string, error) {
	return m.userID, m.valErr
}

func TestAppMiddleware_MissingHeader(t *testing.T) {
	mw := auth.AppMiddleware(&mockAppMiddlewareSvc{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rec.Code)
	}
}

func TestAppMiddleware_NoBearerPrefix(t *testing.T) {
	mw := auth.AppMiddleware(&mockAppMiddlewareSvc{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token abc")
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rec.Code)
	}
}

func TestAppMiddleware_ValidateError(t *testing.T) {
	svc := &mockAppMiddlewareSvc{valErr: errors.New("token revoked")}
	mw := auth.AppMiddleware(svc)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	rec := httptest.NewRecorder()
	mw(okHandler()).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rec.Code)
	}
}

func TestAppMiddleware_ValidToken_SetsUserID(t *testing.T) {
	svc := &mockAppMiddlewareSvc{userID: "user-123"}
	mw := auth.AppMiddleware(svc)

	var capturedID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID, _ = principal.UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	rec := httptest.NewRecorder()
	mw(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rec.Code)
	}
	if capturedID != "user-123" {
		t.Errorf("want user-123, got %q", capturedID)
	}
}

func TestAppMiddleware_RealJWT_Integration(t *testing.T) {
	tokenSvc := auth.NewTokenService("secret", time.Minute, time.Hour)
	raw, err := tokenSvc.GenerateAppToken("user-456", "token-id-1")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	svc := &mockAppMiddlewareSvc{userID: "user-456"}
	mw := auth.AppMiddleware(svc)

	var capturedID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID, _ = principal.UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+raw)
	rec := httptest.NewRecorder()
	mw(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rec.Code)
	}
	if capturedID != "user-456" {
		t.Errorf("want user-456, got %q", capturedID)
	}
}
