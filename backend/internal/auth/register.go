package auth

import (
	"context"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RegisterHandler struct {
	auth AuthService
}

func NewRegisterHandler(auth AuthService) *RegisterHandler {
	return &RegisterHandler{auth: auth}
}

// Method implements handler.Handler.
func (h *RegisterHandler) Method() string              { return "POST" }
// Pattern implements handler.Handler.
func (h *RegisterHandler) Pattern() string             { return "/auth/register" }
// Middlewares implements handler.Handler.
func (h *RegisterHandler) Middlewares() []middleware.Tag { return nil }

// Validate implements handler.Handler.
func (h *RegisterHandler) Validate(req RegisterRequest) error {
	if req.FirstName == "" {
		return errValidation("first_name is required")
	}
	if req.LastName == "" {
		return errValidation("last_name is required")
	}
	if req.Email == "" {
		return errValidation("email is required")
	}
	if req.Password == "" {
		return errValidation("password is required")
	}
	return nil
}

// Serve implements handler.Handler.
func (h *RegisterHandler) Serve(ctx context.Context, req RegisterRequest) (LoginResponse, error) {
	access, refresh, err := h.auth.Register(ctx, req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{Token: access, RefreshToken: refresh}, nil
}
