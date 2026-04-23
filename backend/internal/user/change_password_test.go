package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
	"github.com/basselshurbaji/mr_bean/backend/internal/user"
)

func TestChangePasswordHandler_Validate(t *testing.T) {
	h := user.NewChangePasswordHandler(&mockUserService{})

	if err := h.Validate(user.ChangePasswordRequest{OldPassword: "old", NewPassword: "new"}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if err := h.Validate(user.ChangePasswordRequest{NewPassword: "new"}); err == nil {
		t.Error("expected error when old_password is empty")
	}
	if err := h.Validate(user.ChangePasswordRequest{OldPassword: "old"}); err == nil {
		t.Error("expected error when new_password is empty")
	}
}

func TestChangePasswordHandler_Serve_NoContext(t *testing.T) {
	h := user.NewChangePasswordHandler(&mockUserService{})

	_, err := h.Serve(context.Background(), user.ChangePasswordRequest{OldPassword: "old", NewPassword: "new"})
	if err == nil {
		t.Error("expected error when no user ID in context")
	}
}

func TestChangePasswordHandler_Serve_Success(t *testing.T) {
	h := user.NewChangePasswordHandler(&mockUserService{})

	ctx := principal.WithUserID(context.Background(), "user-123")
	_, err := h.Serve(ctx, user.ChangePasswordRequest{OldPassword: "old", NewPassword: "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChangePasswordHandler_Serve_ServiceError(t *testing.T) {
	svc := &mockUserService{err: errors.New("incorrect password")}
	h := user.NewChangePasswordHandler(svc)

	ctx := principal.WithUserID(context.Background(), "user-123")
	_, err := h.Serve(ctx, user.ChangePasswordRequest{OldPassword: "wrong", NewPassword: "new"})
	if err == nil {
		t.Error("expected error from service")
	}
}
