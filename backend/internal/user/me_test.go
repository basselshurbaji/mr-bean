package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
	"github.com/basselshurbaji/mr_bean/backend/internal/user"
)

type mockUserService struct {
	u   *user.User
	err error
}

func (m *mockUserService) GetByID(_ context.Context, _ string) (*user.User, error) {
	return m.u, m.err
}

func TestMeHandler_Validate(t *testing.T) {
	h := user.NewMeHandler(&mockUserService{})
	if err := h.Validate(user.MeRequest{}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestMeHandler_Serve_NoContext(t *testing.T) {
	h := user.NewMeHandler(&mockUserService{})

	if _, err := h.Serve(context.Background(), user.MeRequest{}); err == nil {
		t.Error("expected error when no user ID in context")
	}
}

func TestMeHandler_Serve_ReturnsUser(t *testing.T) {
	now := time.Now()
	svc := &mockUserService{
		u: &user.User{
			ID:        "user-123",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			CreatedAt: now,
		},
	}
	h := user.NewMeHandler(svc)

	ctx := principal.WithUserID(context.Background(), "user-123")
	res, err := h.Serve(ctx, user.MeRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "user-123" {
		t.Errorf("expected ID user-123, got %s", res.ID)
	}
	if res.Email != "john@doe.com" {
		t.Errorf("expected email john@doe.com, got %s", res.Email)
	}
}

func TestMeHandler_Serve_ServiceError(t *testing.T) {
	svc := &mockUserService{err: errors.New("db error")}
	h := user.NewMeHandler(svc)

	ctx := principal.WithUserID(context.Background(), "user-123")
	if _, err := h.Serve(ctx, user.MeRequest{}); err == nil {
		t.Error("expected error from service")
	}
}
