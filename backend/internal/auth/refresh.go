package auth

import (
	"context"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshHandler struct {
	auth AuthService
}

func NewRefreshHandler(auth AuthService) *RefreshHandler {
	return &RefreshHandler{auth: auth}
}

func (h *RefreshHandler) Method() string              { return "POST" }
func (h *RefreshHandler) Pattern() string             { return "/auth/refresh" }
func (h *RefreshHandler) Middlewares() []middleware.Tag { return nil }

func (h *RefreshHandler) Validate(req RefreshRequest) error {
	if req.RefreshToken == "" {
		return errValidation("refresh_token is required")
	}
	return nil
}

func (h *RefreshHandler) Serve(ctx context.Context, req RefreshRequest) (RefreshResponse, error) {
	access, refresh, err := h.auth.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return RefreshResponse{}, err
	}
	return RefreshResponse{Token: access, RefreshToken: refresh}, nil
}
