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
