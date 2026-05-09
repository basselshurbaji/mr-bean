package auth_test

import (
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/auth"
)

func newTestTokenService() auth.TokenService {
	return auth.NewTokenService("test-secret", time.Minute, time.Hour)
}

func TestTokenService_AccessRoundTrip(t *testing.T) {
	svc := newTestTokenService()

	token, err := svc.GenerateAccessToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	claims, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("expected user-123, got %s", claims.UserID)
	}
	if claims.TokenType != auth.TokenTypeAccess {
		t.Errorf("expected access type, got %s", claims.TokenType)
	}
}

func TestTokenService_RefreshRoundTrip(t *testing.T) {
	svc := newTestTokenService()

	token, err := svc.GenerateRefreshToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	claims, err := svc.ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("expected user-123, got %s", claims.UserID)
	}
	if claims.TokenType != auth.TokenTypeRefresh {
		t.Errorf("expected refresh type, got %s", claims.TokenType)
	}
}

func TestTokenService_AccessRejectsRefreshToken(t *testing.T) {
	svc := newTestTokenService()

	refresh, err := svc.GenerateRefreshToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	if _, err := svc.ValidateAccessToken(refresh); err == nil {
		t.Error("expected error when validating refresh token as access token")
	}
}

func TestTokenService_RefreshRejectsAccessToken(t *testing.T) {
	svc := newTestTokenService()

	access, err := svc.GenerateAccessToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	if _, err := svc.ValidateRefreshToken(access); err == nil {
		t.Error("expected error when validating access token as refresh token")
	}
}

func TestTokenService_ExpiredToken(t *testing.T) {
	svc := auth.NewTokenService("test-secret", -time.Second, time.Hour)

	token, err := svc.GenerateAccessToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	if _, err := svc.ValidateAccessToken(token); err == nil {
		t.Error("expected error for expired token")
	}
}

func TestTokenService_AppRoundTrip(t *testing.T) {
	svc := newTestTokenService()

	raw, err := svc.GenerateAppToken("user-123", "token-id-abc")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	claims, err := svc.ValidateAppToken(raw)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("expected user-123, got %s", claims.UserID)
	}
	if claims.TokenType != auth.TokenTypeApp {
		t.Errorf("expected app type, got %s", claims.TokenType)
	}
	if claims.ID != "token-id-abc" {
		t.Errorf("expected jti token-id-abc, got %s", claims.ID)
	}
	if claims.ExpiresAt != nil {
		t.Error("expected no expiry on app token")
	}
}

func TestTokenService_AppRejectsAccessToken(t *testing.T) {
	svc := newTestTokenService()

	access, err := svc.GenerateAccessToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	if _, err := svc.ValidateAppToken(access); err == nil {
		t.Error("expected error when validating access token as app token")
	}
}

func TestTokenService_AppRejectsRefreshToken(t *testing.T) {
	svc := newTestTokenService()

	refresh, err := svc.GenerateRefreshToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	if _, err := svc.ValidateAppToken(refresh); err == nil {
		t.Error("expected error when validating refresh token as app token")
	}
}

func TestTokenService_AccessRejectsAppToken(t *testing.T) {
	svc := newTestTokenService()

	app, err := svc.GenerateAppToken("user-123", "token-id-abc")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	if _, err := svc.ValidateAccessToken(app); err == nil {
		t.Error("expected error when validating app token as access token")
	}
}
