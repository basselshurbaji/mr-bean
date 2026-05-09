package auth

import (
	"context"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type CreateAppTokenRequest struct {
	AppName string `json:"app_name"`
}

type CreateAppTokenResponse struct {
	ID        string    `json:"id"`
	AppName   string    `json:"app_name"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAppTokenHandler struct {
	svc AppTokenService
}

func NewCreateAppTokenHandler(svc AppTokenService) *CreateAppTokenHandler {
	return &CreateAppTokenHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *CreateAppTokenHandler) Method() string { return "POST" }

// Pattern implements handler.Handler.
func (h *CreateAppTokenHandler) Pattern() string { return "/app-token" }

// Middlewares implements handler.Handler.
func (h *CreateAppTokenHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

// Validate implements handler.Handler.
func (h *CreateAppTokenHandler) Validate(req CreateAppTokenRequest) error {
	if req.AppName == "" {
		return errValidation("app_name is required")
	}
	return nil
}

// Serve implements handler.Handler.
func (h *CreateAppTokenHandler) Serve(ctx context.Context, req CreateAppTokenRequest) (CreateAppTokenResponse, error) {
	userID, _ := principal.UserIDFromContext(ctx)
	token, record, err := h.svc.Create(ctx, userID, req.AppName)
	if err != nil {
		return CreateAppTokenResponse{}, err
	}
	return CreateAppTokenResponse{
		ID:        record.ID,
		AppName:   record.AppName,
		Token:     token,
		CreatedAt: record.CreatedAt,
	}, nil
}