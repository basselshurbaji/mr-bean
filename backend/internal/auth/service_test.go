package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/basselshurbaji/mr_bean/backend/internal/auth"
	"github.com/basselshurbaji/mr_bean/backend/internal/mailer"
	"github.com/basselshurbaji/mr_bean/backend/internal/user"
)

type mockMailer struct{}

func (m *mockMailer) Send(_ context.Context, _ mailer.Email) error {
	return nil
}

type mockUserStore struct {
	user        *user.User
	err         error
	createdUser *user.User
	createErr   error
}

func (m *mockUserStore) GetByEmail(_ context.Context, _ string) (*user.User, error) {
	return m.user, m.err
}

func (m *mockUserStore) GetByID(_ context.Context, _ string) (*user.User, error) {
	return nil, nil
}

func (m *mockUserStore) Create(_ context.Context, _, _, _, _ string) (*user.User, error) {
	return m.createdUser, m.createErr
}

func (m *mockUserStore) UpdateProfile(_ context.Context, _, _, _ string) (*user.User, error) {
	return nil, nil
}

func (m *mockUserStore) UpdatePassword(_ context.Context, _, _ string) error {
	return nil
}

func hashedPassword(t *testing.T, plain string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	return string(h)
}

func newTestAuthService(store user.UserRepo) auth.AuthService {
	tokens := auth.NewTokenService("test-secret", time.Minute, time.Hour)
	return auth.NewAuthService(store, tokens, &mockMailer{})
}

func TestAuthService_Login_HappyPath(t *testing.T) {
	store := &mockUserStore{
		user: &user.User{
			ID:           "user-123",
			PasswordHash: hashedPassword(t, "correct-password"),
			IsActive:     true,
		},
	}
	svc := newTestAuthService(store)

	access, refresh, err := svc.Login(context.Background(), "a@b.com", "correct-password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access == "" || refresh == "" {
		t.Error("expected non-empty tokens")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	store := &mockUserStore{err: errors.New("not found")}
	svc := newTestAuthService(store)

	if _, _, err := svc.Login(context.Background(), "a@b.com", "pw"); err == nil {
		t.Error("expected error for missing user")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	store := &mockUserStore{
		user: &user.User{
			ID:           "user-123",
			PasswordHash: hashedPassword(t, "correct-password"),
			IsActive:     true,
		},
	}
	svc := newTestAuthService(store)

	if _, _, err := svc.Login(context.Background(), "a@b.com", "wrong-password"); err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	store := &mockUserStore{
		user: &user.User{
			ID:           "user-123",
			PasswordHash: hashedPassword(t, "password"),
			IsActive:     false,
		},
	}
	svc := newTestAuthService(store)

	if _, _, err := svc.Login(context.Background(), "a@b.com", "password"); err == nil {
		t.Error("expected error for inactive user")
	}
}

func TestAuthService_Refresh_HappyPath(t *testing.T) {
	store := &mockUserStore{}
	tokens := auth.NewTokenService("test-secret", time.Minute, time.Hour)
	svc := auth.NewAuthService(store, tokens, &mockMailer{})

	refreshToken, err := tokens.GenerateRefreshToken("user-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	access, refresh, err := svc.Refresh(context.Background(), refreshToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access == "" || refresh == "" {
		t.Error("expected non-empty tokens")
	}
}

func TestAuthService_Refresh_InvalidToken(t *testing.T) {
	svc := newTestAuthService(&mockUserStore{})

	if _, _, err := svc.Refresh(context.Background(), "not-a-token"); err == nil {
		t.Error("expected error for invalid refresh token")
	}
}

func TestAuthService_Register_HappyPath(t *testing.T) {
	store := &mockUserStore{
		createdUser: &user.User{ID: "new-user-123", IsActive: true},
	}
	svc := newTestAuthService(store)

	access, refresh, err := svc.Register(context.Background(), "Jane", "Doe", "jane@example.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access == "" || refresh == "" {
		t.Error("expected non-empty tokens")
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	store := &mockUserStore{createErr: errors.New("unique constraint violation")}
	svc := newTestAuthService(store)

	if _, _, err := svc.Register(context.Background(), "Jane", "Doe", "existing@example.com", "password123"); err == nil {
		t.Error("expected error for duplicate email")
	}
}
