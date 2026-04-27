package extraction

import (
	"context"
	"errors"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type ListExtractionsRequest struct {
	Limit int `schema:"limit"`
	Page  int `schema:"page"`
}

type ListExtractionsHandler struct {
	svc ExtractionService
}

func NewListExtractionsHandler(svc ExtractionService) *ListExtractionsHandler {
	return &ListExtractionsHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *ListExtractionsHandler) Method() string { return "GET" }

// Pattern implements handler.Handler.
func (h *ListExtractionsHandler) Pattern() string { return "/extractions" }

// Middlewares implements handler.Handler.
func (h *ListExtractionsHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

// Validate implements handler.Handler.
func (h *ListExtractionsHandler) Validate(_ ListExtractionsRequest) error { return nil }

// Serve implements handler.Handler.
func (h *ListExtractionsHandler) Serve(ctx context.Context, req ListExtractionsRequest) ([]ExtractionResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}

	items, err := h.svc.ListExtractions(ctx, userID, limit, page)
	if err != nil {
		return nil, err
	}
	res := make([]ExtractionResponse, len(items))
	for i, e := range items {
		res[i] = toResponse(e)
	}
	return res, nil
}
