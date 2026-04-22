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
