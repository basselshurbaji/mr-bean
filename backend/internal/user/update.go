package user

import (
	"context"
	"errors"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type UpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateHandler struct {
	users UserService
}

func NewUpdateHandler(users UserService) *UpdateHandler {
	return &UpdateHandler{users: users}
}

func (h *UpdateHandler) Method() string  { return "PATCH" }
func (h *UpdateHandler) Pattern() string { return "/user/me" }

func (h *UpdateHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

func (h *UpdateHandler) Validate(req UpdateRequest) error {
	if req.FirstName == "" && req.LastName == "" {
		return errors.New("at least one of first_name or last_name is required")
	}
	return nil
}

func (h *UpdateHandler) Serve(ctx context.Context, req UpdateRequest) (UpdateResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return UpdateResponse{}, errors.New("unauthorized")
	}

	u, err := h.users.UpdateProfile(ctx, userID, req.FirstName, req.LastName)
	if err != nil {
		return UpdateResponse{}, err
	}

	return UpdateResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}
