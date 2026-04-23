package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/basselshurbaji/mr_bean/backend/internal/auth"
)

type mockAuthService struct {
	accessToken  string
	refreshToken string
	err          error
}

func (m *mockAuthService) Login(_ context.Context, _, _ string) (string, string, error) {
	return m.accessToken, m.refreshToken, m.err
}

func (m *mockAuthService) Refresh(_ context.Context, _ string) (string, string, error) {
	return m.accessToken, m.refreshToken, m.err
}

func (m *mockAuthService) Register(_ context.Context, _, _, _, _ string) (string, string, error) {
	return m.accessToken, m.refreshToken, m.err
}

// --- LoginHandler ---

func TestLoginHandler_Validate_EmptyEmail(t *testing.T) {
	h := auth.NewLoginHandler(&mockAuthService{})
	if err := h.Validate(auth.LoginRequest{Password: "pw"}); err == nil {
		t.Error("expected error for empty email")
	}
}

func TestLoginHandler_Validate_EmptyPassword(t *testing.T) {
	h := auth.NewLoginHandler(&mockAuthService{})
	if err := h.Validate(auth.LoginRequest{Email: "a@b.com"}); err == nil {
		t.Error("expected error for empty password")
	}
}

func TestLoginHandler_Validate_Valid(t *testing.T) {
	h := auth.NewLoginHandler(&mockAuthService{})
	if err := h.Validate(auth.LoginRequest{Email: "a@b.com", Password: "pw"}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestLoginHandler_Serve_ReturnsTokens(t *testing.T) {
	svc := &mockAuthService{accessToken: "access", refreshToken: "refresh"}
	h := auth.NewLoginHandler(svc)

	res, err := h.Serve(context.Background(), auth.LoginRequest{Email: "a@b.com", Password: "pw"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Token != "access" || res.RefreshToken != "refresh" {
		t.Errorf("unexpected tokens: %+v", res)
	}
}

func TestLoginHandler_Serve_PropagatesError(t *testing.T) {
	svc := &mockAuthService{err: errors.New("invalid credentials")}
	h := auth.NewLoginHandler(svc)

	if _, err := h.Serve(context.Background(), auth.LoginRequest{Email: "a@b.com", Password: "wrong"}); err == nil {
		t.Error("expected error from service")
	}
}

// --- RefreshHandler ---

func TestRefreshHandler_Validate_EmptyToken(t *testing.T) {
	h := auth.NewRefreshHandler(&mockAuthService{})
	if err := h.Validate(auth.RefreshRequest{}); err == nil {
		t.Error("expected error for empty refresh_token")
	}
}

func TestRefreshHandler_Validate_Valid(t *testing.T) {
	h := auth.NewRefreshHandler(&mockAuthService{})
	if err := h.Validate(auth.RefreshRequest{RefreshToken: "token"}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRefreshHandler_Serve_ReturnsTokens(t *testing.T) {
	svc := &mockAuthService{accessToken: "new-access", refreshToken: "new-refresh"}
	h := auth.NewRefreshHandler(svc)

	res, err := h.Serve(context.Background(), auth.RefreshRequest{RefreshToken: "old-refresh"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Token != "new-access" || res.RefreshToken != "new-refresh" {
		t.Errorf("unexpected tokens: %+v", res)
	}
}

func TestRefreshHandler_Serve_PropagatesError(t *testing.T) {
	svc := &mockAuthService{err: errors.New("invalid refresh token")}
	h := auth.NewRefreshHandler(svc)

	if _, err := h.Serve(context.Background(), auth.RefreshRequest{RefreshToken: "bad"}); err == nil {
		t.Error("expected error from service")
	}
}

// --- RegisterHandler ---

func TestRegisterHandler_Validate_MissingFirstName(t *testing.T) {
	h := auth.NewRegisterHandler(&mockAuthService{})
	if err := h.Validate(auth.RegisterRequest{LastName: "Doe", Email: "a@b.com", Password: "pw"}); err == nil {
		t.Error("expected error for missing first_name")
	}
}

func TestRegisterHandler_Validate_MissingLastName(t *testing.T) {
	h := auth.NewRegisterHandler(&mockAuthService{})
	if err := h.Validate(auth.RegisterRequest{FirstName: "Jane", Email: "a@b.com", Password: "pw"}); err == nil {
		t.Error("expected error for missing last_name")
	}
}

func TestRegisterHandler_Validate_MissingEmail(t *testing.T) {
	h := auth.NewRegisterHandler(&mockAuthService{})
	if err := h.Validate(auth.RegisterRequest{FirstName: "Jane", LastName: "Doe", Password: "pw"}); err == nil {
		t.Error("expected error for missing email")
	}
}

func TestRegisterHandler_Validate_MissingPassword(t *testing.T) {
	h := auth.NewRegisterHandler(&mockAuthService{})
	if err := h.Validate(auth.RegisterRequest{FirstName: "Jane", LastName: "Doe", Email: "a@b.com"}); err == nil {
		t.Error("expected error for missing password")
	}
}

func TestRegisterHandler_Validate_Valid(t *testing.T) {
	h := auth.NewRegisterHandler(&mockAuthService{})
	if err := h.Validate(auth.RegisterRequest{FirstName: "Jane", LastName: "Doe", Email: "a@b.com", Password: "pw"}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRegisterHandler_Serve_ReturnsTokens(t *testing.T) {
	svc := &mockAuthService{accessToken: "access", refreshToken: "refresh"}
	h := auth.NewRegisterHandler(svc)

	res, err := h.Serve(context.Background(), auth.RegisterRequest{
		FirstName: "Jane", LastName: "Doe", Email: "a@b.com", Password: "pw",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Token != "access" || res.RefreshToken != "refresh" {
		t.Errorf("unexpected tokens: %+v", res)
	}
}

func TestRegisterHandler_Serve_PropagatesError(t *testing.T) {
	svc := &mockAuthService{err: errors.New("email already registered")}
	h := auth.NewRegisterHandler(svc)

	if _, err := h.Serve(context.Background(), auth.RegisterRequest{
		FirstName: "Jane", LastName: "Doe", Email: "existing@b.com", Password: "pw",
	}); err == nil {
		t.Error("expected error from service")
	}
}
