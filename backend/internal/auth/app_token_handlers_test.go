package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/auth"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type mockAppTokenService struct {
	token     string
	record    *auth.AppToken
	createErr error
	revokeErr error
	userID    string
	valErr    error
}

func (m *mockAppTokenService) Create(_ context.Context, _, _ string) (string, *auth.AppToken, error) {
	return m.token, m.record, m.createErr
}

func (m *mockAppTokenService) Revoke(_ context.Context, _, _ string) error {
	return m.revokeErr
}

func (m *mockAppTokenService) Validate(_ context.Context, _ string) (string, error) {
	return m.userID, m.valErr
}

// --- CreateAppTokenHandler ---

func TestCreateAppTokenHandler_Validate_MissingAppName(t *testing.T) {
	h := auth.NewCreateAppTokenHandler(&mockAppTokenService{})
	if err := h.Validate(auth.CreateAppTokenRequest{}); err == nil {
		t.Error("expected error for missing app_name")
	}
}

func TestCreateAppTokenHandler_Validate_Valid(t *testing.T) {
	h := auth.NewCreateAppTokenHandler(&mockAppTokenService{})
	if err := h.Validate(auth.CreateAppTokenRequest{AppName: "my-app"}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateAppTokenHandler_Serve_HappyPath(t *testing.T) {
	record := &auth.AppToken{ID: "tok-1", AppName: "my-app", CreatedAt: time.Now()}
	svc := &mockAppTokenService{token: "signed.jwt", record: record}
	h := auth.NewCreateAppTokenHandler(svc)

	ctx := principal.WithUserID(context.Background(), "user-123")
	res, err := h.Serve(ctx, auth.CreateAppTokenRequest{AppName: "my-app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "tok-1" {
		t.Errorf("expected tok-1, got %s", res.ID)
	}
	if res.Token != "signed.jwt" {
		t.Errorf("expected signed.jwt, got %s", res.Token)
	}
	if res.AppName != "my-app" {
		t.Errorf("expected my-app, got %s", res.AppName)
	}
}

func TestCreateAppTokenHandler_Serve_PropagatesError(t *testing.T) {
	svc := &mockAppTokenService{createErr: errors.New("db error")}
	h := auth.NewCreateAppTokenHandler(svc)

	ctx := principal.WithUserID(context.Background(), "user-123")
	if _, err := h.Serve(ctx, auth.CreateAppTokenRequest{AppName: "my-app"}); err == nil {
		t.Error("expected error from service")
	}
}

// --- RevokeAppTokenHandler ---

func TestRevokeAppTokenHandler_Validate_MissingID(t *testing.T) {
	h := auth.NewRevokeAppTokenHandler(&mockAppTokenService{})
	if err := h.Validate(auth.RevokeAppTokenRequest{}); err == nil {
		t.Error("expected error for missing id")
	}
}

func TestRevokeAppTokenHandler_Validate_Valid(t *testing.T) {
	h := auth.NewRevokeAppTokenHandler(&mockAppTokenService{})
	if err := h.Validate(auth.RevokeAppTokenRequest{ID: "tok-1"}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRevokeAppTokenHandler_Serve_HappyPath(t *testing.T) {
	svc := &mockAppTokenService{}
	h := auth.NewRevokeAppTokenHandler(svc)

	ctx := principal.WithUserID(context.Background(), "user-123")
	if _, err := h.Serve(ctx, auth.RevokeAppTokenRequest{ID: "tok-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRevokeAppTokenHandler_Serve_PropagatesError(t *testing.T) {
	svc := &mockAppTokenService{revokeErr: errors.New("not found")}
	h := auth.NewRevokeAppTokenHandler(svc)

	ctx := principal.WithUserID(context.Background(), "user-123")
	if _, err := h.Serve(ctx, auth.RevokeAppTokenRequest{ID: "tok-1"}); err == nil {
		t.Error("expected error from service")
	}
}
