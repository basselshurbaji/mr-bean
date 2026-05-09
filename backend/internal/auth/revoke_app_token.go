package auth

import (
	"context"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type RevokeAppTokenRequest struct {
	ID string `url:"id"`
}

type RevokeAppTokenHandler struct {
	svc AppTokenService
}

func NewRevokeAppTokenHandler(svc AppTokenService) *RevokeAppTokenHandler {
	return &RevokeAppTokenHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *RevokeAppTokenHandler) Method() string { return "DELETE" }

// Pattern implements handler.Handler.
func (h *RevokeAppTokenHandler) Pattern() string { return "/app-token/{id}" }

// Middlewares implements handler.Handler.
func (h *RevokeAppTokenHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

// Validate implements handler.Handler.
func (h *RevokeAppTokenHandler) Validate(req RevokeAppTokenRequest) error {
	if req.ID == "" {
		return errValidation("id is required")
	}
	return nil
}

// Serve implements handler.Handler.
func (h *RevokeAppTokenHandler) Serve(ctx context.Context, req RevokeAppTokenRequest) (handler.NoContent, error) {
	userID, _ := principal.UserIDFromContext(ctx)
	if err := h.svc.Revoke(ctx, req.ID, userID); err != nil {
		return handler.NoContent{}, err
	}
	return handler.NoContent{}, nil
}