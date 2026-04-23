package user

import (
	"context"
	"errors"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct{}

type ChangePasswordHandler struct {
	users UserService
}

func NewChangePasswordHandler(users UserService) *ChangePasswordHandler {
	return &ChangePasswordHandler{users: users}
}

func (h *ChangePasswordHandler) Method() string  { return "POST" }
func (h *ChangePasswordHandler) Pattern() string { return "/user/change-password" }

func (h *ChangePasswordHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

func (h *ChangePasswordHandler) Validate(req ChangePasswordRequest) error {
	if req.OldPassword == "" {
		return errors.New("old_password is required")
	}
	if req.NewPassword == "" {
		return errors.New("new_password is required")
	}
	return nil
}

func (h *ChangePasswordHandler) Serve(ctx context.Context, req ChangePasswordRequest) (ChangePasswordResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return ChangePasswordResponse{}, errors.New("unauthorized")
	}

	if err := h.users.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword); err != nil {
		return ChangePasswordResponse{}, err
	}

	return ChangePasswordResponse{}, nil
}
