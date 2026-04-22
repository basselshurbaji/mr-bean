package user

import (
	"context"
	"errors"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type MeRequest struct{}

type MeResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type MeHandler struct {
	users UserService
}

func NewMeHandler(users UserService) *MeHandler {
	return &MeHandler{users: users}
}

func (h *MeHandler) Method() string  { return "GET" }
func (h *MeHandler) Pattern() string { return "/user/me" }

func (h *MeHandler) Validate(_ MeRequest) error { return nil }

func (h *MeHandler) Serve(ctx context.Context, _ MeRequest) (MeResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return MeResponse{}, errors.New("unauthorized")
	}

	u, err := h.users.GetByID(ctx, userID)
	if err != nil {
		return MeResponse{}, err
	}

	return MeResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}
