package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/user"
)

type mockUserRepo struct {
	u   *user.User
	err error
}

func (m *mockUserRepo) GetByEmail(_ context.Context, _ string) (*user.User, error) {
	return m.u, m.err
}

func (m *mockUserRepo) GetByID(_ context.Context, _ string) (*user.User, error) {
	return m.u, m.err
}

func (m *mockUserRepo) Create(_ context.Context, _, _, _, _ string) (*user.User, error) {
	return nil, nil
}

func (m *mockUserRepo) UpdateProfile(_ context.Context, _, _, _ string) (*user.User, error) {
	return m.u, m.err
}

func (m *mockUserRepo) UpdatePassword(_ context.Context, _, _ string) error {
	return m.err
}

func TestUserService_GetByID_Found(t *testing.T) {
	expected := &user.User{ID: "user-123", Email: "a@b.com", CreatedAt: time.Now()}
	svc := user.NewUserService(&mockUserRepo{u: expected})

	got, err := svc.GetByID(context.Background(), "user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != expected.ID {
		t.Errorf("expected ID %s, got %s", expected.ID, got.ID)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	svc := user.NewUserService(&mockUserRepo{err: errors.New("not found")})

	if _, err := svc.GetByID(context.Background(), "missing"); err == nil {
		t.Error("expected error for missing user")
	}
}

func TestUserService_UpdateProfile_BothFields(t *testing.T) {
	expected := &user.User{ID: "user-123", FirstName: "Jane", LastName: "Doe"}
	svc := user.NewUserService(&mockUserRepo{u: expected})

	got, err := svc.UpdateProfile(context.Background(), "user-123", "Jane", "Doe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.FirstName != "Jane" || got.LastName != "Doe" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestUserService_UpdateProfile_PartialFillsFromExisting(t *testing.T) {
	existing := &user.User{ID: "user-123", FirstName: "John", LastName: "Doe"}
	repo := &mockUserRepo{u: existing}
	svc := user.NewUserService(repo)

	_, err := svc.UpdateProfile(context.Background(), "user-123", "Jane", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserService_ChangePassword_IncorrectOld(t *testing.T) {
	u := &user.User{ID: "user-123", PasswordHash: "$2a$10$invalidhash"}
	svc := user.NewUserService(&mockUserRepo{u: u})

	err := svc.ChangePassword(context.Background(), "user-123", "wrong", "newpass")
	if err == nil {
		t.Error("expected error for incorrect old password")
	}
}
