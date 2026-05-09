package bean

import (
	"context"
	"errors"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type ListBeansRequest struct{}

type ListBeansHandler struct {
	svc BeanService
}

func NewListBeansHandler(svc BeanService) *ListBeansHandler {
	return &ListBeansHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *ListBeansHandler) Method() string  { return "GET" }
// Pattern implements handler.Handler.
func (h *ListBeansHandler) Pattern() string { return "/beans" }

// Middlewares implements handler.Handler.
func (h *ListBeansHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAnyAuthenticated}
}

// Validate implements handler.Handler.
func (h *ListBeansHandler) Validate(_ ListBeansRequest) error { return nil }

// Serve implements handler.Handler.
func (h *ListBeansHandler) Serve(ctx context.Context, _ ListBeansRequest) ([]BeanResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	items, err := h.svc.ListBeans(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make([]BeanResponse, len(items))
	for i, b := range items {
		res[i] = beanToResponse(b)
	}
	return res, nil
}
