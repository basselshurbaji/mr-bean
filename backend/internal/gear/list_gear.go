package gear

import (
	"context"
	"errors"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type ListGearRequest struct{}

type ListGearHandler struct {
	svc GearService
}

func NewListGearHandler(svc GearService) *ListGearHandler {
	return &ListGearHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *ListGearHandler) Method() string  { return "GET" }
// Pattern implements handler.Handler.
func (h *ListGearHandler) Pattern() string { return "/gear" }

// Middlewares implements handler.Handler.
func (h *ListGearHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

// Validate implements handler.Handler.
func (h *ListGearHandler) Validate(_ ListGearRequest) error { return nil }

// Serve implements handler.Handler.
func (h *ListGearHandler) Serve(ctx context.Context, _ ListGearRequest) ([]GearResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	items, err := h.svc.ListGear(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make([]GearResponse, len(items))
	for i, g := range items {
		res[i] = gearToResponse(g)
	}
	return res, nil
}
