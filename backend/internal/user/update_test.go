package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
	"github.com/basselshurbaji/mr_bean/backend/internal/user"
)

func TestUpdateHandler_Validate(t *testing.T) {
	h := user.NewUpdateHandler(&mockUserService{})

	if err := h.Validate(user.UpdateRequest{FirstName: "Jane"}); err != nil {
		t.Errorf("expected no error with first_name only, got %v", err)
	}
	if err := h.Validate(user.UpdateRequest{LastName: "Doe"}); err != nil {
		t.Errorf("expected no error with last_name only, got %v", err)
	}
	if err := h.Validate(user.UpdateRequest{}); err == nil {
		t.Error("expected error when both fields are empty")
	}
}

func TestUpdateHandler_Serve_NoContext(t *testing.T) {
	h := user.NewUpdateHandler(&mockUserService{})

	_, err := h.Serve(context.Background(), user.UpdateRequest{FirstName: "Jane"})
	if err == nil {
		t.Error("expected error when no user ID in context")
	}
}

func TestUpdateHandler_Serve_ReturnsUpdatedUser(t *testing.T) {
	now := time.Now()
	svc := &mockUserService{
		u: &user.User{
			ID:        "user-123",
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane@doe.com",
			CreatedAt: now,
		},
	}
	h := user.NewUpdateHandler(svc)

	ctx := principal.WithUserID(context.Background(), "user-123")
	res, err := h.Serve(ctx, user.UpdateRequest{FirstName: "Jane", LastName: "Doe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FirstName != "Jane" {
		t.Errorf("expected first_name Jane, got %s", res.FirstName)
	}
	if res.LastName != "Doe" {
		t.Errorf("expected last_name Doe, got %s", res.LastName)
	}
}

func TestUpdateHandler_Serve_ServiceError(t *testing.T) {
	svc := &mockUserService{err: errors.New("db error")}
	h := user.NewUpdateHandler(svc)

	ctx := principal.WithUserID(context.Background(), "user-123")
	_, err := h.Serve(ctx, user.UpdateRequest{FirstName: "Jane"})
	if err == nil {
		t.Error("expected error from service")
	}
}
