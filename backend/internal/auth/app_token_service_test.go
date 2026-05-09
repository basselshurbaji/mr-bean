package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/auth"
)

// --- mocks ---

type mockAppTokenRepo struct {
	record    *auth.AppToken
	createErr error
	getErr    error
	revokeErr error
}

func (m *mockAppTokenRepo) Create(_ context.Context, _, _ string) (*auth.AppToken, error) {
	return m.record, m.createErr
}

func (m *mockAppTokenRepo) GetByID(_ context.Context, _ string) (*auth.AppToken, error) {
	return m.record, m.getErr
}

func (m *mockAppTokenRepo) Revoke(_ context.Context, _, _ string) error {
	return m.revokeErr
}

type mockTokenSvc struct {
	appToken string
	genErr   error
	claims   *auth.Claims
	valErr   error
}

func (m *mockTokenSvc) GenerateAccessToken(_ string) (string, error)          { return "", nil }
func (m *mockTokenSvc) GenerateRefreshToken(_ string) (string, error)         { return "", nil }
func (m *mockTokenSvc) GenerateAppToken(_, _ string) (string, error)          { return m.appToken, m.genErr }
func (m *mockTokenSvc) ValidateAccessToken(_ string) (*auth.Claims, error)    { return nil, nil }
func (m *mockTokenSvc) ValidateRefreshToken(_ string) (*auth.Claims, error)   { return nil, nil }
func (m *mockTokenSvc) ValidateAppToken(_ string) (*auth.Claims, error)       { return m.claims, m.valErr }

// --- AppTokenService.Create ---

func TestAppTokenService_Create_HappyPath(t *testing.T) {
	record := &auth.AppToken{ID: "tok-1", UserID: "user-123", AppName: "my-app", CreatedAt: time.Now()}
	repo := &mockAppTokenRepo{record: record}
	tokens := &mockTokenSvc{appToken: "signed.jwt.here"}
	svc := auth.NewAppTokenService(repo, tokens)

	token, got, err := svc.Create(context.Background(), "user-123", "my-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "signed.jwt.here" {
		t.Errorf("expected signed.jwt.here, got %s", token)
	}
	if got.ID != record.ID || got.AppName != record.AppName {
		t.Errorf("unexpected record: %+v", got)
	}
}

func TestAppTokenService_Create_RepoError(t *testing.T) {
	repo := &mockAppTokenRepo{createErr: errors.New("db error")}
	svc := auth.NewAppTokenService(repo, &mockTokenSvc{})

	if _, _, err := svc.Create(context.Background(), "user-123", "my-app"); err == nil {
		t.Error("expected error from repo")
	}
}

func TestAppTokenService_Create_TokenGenError(t *testing.T) {
	record := &auth.AppToken{ID: "tok-1", UserID: "user-123", AppName: "my-app"}
	repo := &mockAppTokenRepo{record: record}
	tokens := &mockTokenSvc{genErr: errors.New("signing failed")}
	svc := auth.NewAppTokenService(repo, tokens)

	if _, _, err := svc.Create(context.Background(), "user-123", "my-app"); err == nil {
		t.Error("expected error from token generation")
	}
}

// --- AppTokenService.Revoke ---

func TestAppTokenService_Revoke_HappyPath(t *testing.T) {
	repo := &mockAppTokenRepo{}
	svc := auth.NewAppTokenService(repo, &mockTokenSvc{})

	if err := svc.Revoke(context.Background(), "tok-1", "user-123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAppTokenService_Revoke_RepoError(t *testing.T) {
	repo := &mockAppTokenRepo{revokeErr: errors.New("not found")}
	svc := auth.NewAppTokenService(repo, &mockTokenSvc{})

	if err := svc.Revoke(context.Background(), "tok-1", "user-123"); err == nil {
		t.Error("expected error from repo")
	}
}

// --- AppTokenService.Validate ---

func TestAppTokenService_Validate_HappyPath(t *testing.T) {
	claims := &auth.Claims{}
	claims.ID = "tok-1"
	claims.UserID = "user-123"
	record := &auth.AppToken{ID: "tok-1", UserID: "user-123", Revoked: false}
	repo := &mockAppTokenRepo{record: record}
	tokens := &mockTokenSvc{claims: claims}
	svc := auth.NewAppTokenService(repo, tokens)

	userID, err := svc.Validate(context.Background(), "some.raw.token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "user-123" {
		t.Errorf("expected user-123, got %s", userID)
	}
}

func TestAppTokenService_Validate_InvalidJWT(t *testing.T) {
	tokens := &mockTokenSvc{valErr: errors.New("invalid token")}
	svc := auth.NewAppTokenService(&mockAppTokenRepo{}, tokens)

	if _, err := svc.Validate(context.Background(), "bad.token"); err == nil {
		t.Error("expected error for invalid JWT")
	}
}

func TestAppTokenService_Validate_TokenNotFound(t *testing.T) {
	claims := &auth.Claims{}
	claims.ID = "tok-missing"
	tokens := &mockTokenSvc{claims: claims}
	repo := &mockAppTokenRepo{getErr: errors.New("sql: no rows")}
	svc := auth.NewAppTokenService(repo, tokens)

	if _, err := svc.Validate(context.Background(), "some.raw.token"); err == nil {
		t.Error("expected error for missing DB record")
	}
}

func TestAppTokenService_Validate_TokenRevoked(t *testing.T) {
	claims := &auth.Claims{}
	claims.ID = "tok-1"
	record := &auth.AppToken{ID: "tok-1", Revoked: true}
	tokens := &mockTokenSvc{claims: claims}
	repo := &mockAppTokenRepo{record: record}
	svc := auth.NewAppTokenService(repo, tokens)

	if _, err := svc.Validate(context.Background(), "some.raw.token"); err == nil {
		t.Error("expected error for revoked token")
	}
}
