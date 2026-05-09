package extraction

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type GetExtractionRequest struct {
	ID string `url:"id"`
}

type GetExtractionHandler struct {
	svc ExtractionService
}

func NewGetExtractionHandler(svc ExtractionService) *GetExtractionHandler {
	return &GetExtractionHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *GetExtractionHandler) Method() string { return "GET" }

// Pattern implements handler.Handler.
func (h *GetExtractionHandler) Pattern() string { return "/extractions/{id}" }

// Middlewares implements handler.Handler.
func (h *GetExtractionHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAnyAuthenticated}
}

// Validate implements handler.Handler.
func (h *GetExtractionHandler) Validate(req GetExtractionRequest) error {
	if req.ID == "" {
		return errors.New("id is required")
	}
	return nil
}

// Serve implements handler.Handler.
func (h *GetExtractionHandler) Serve(ctx context.Context, req GetExtractionRequest) (ExtractionResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return ExtractionResponse{}, errors.New("unauthorized")
	}
	e, err := h.svc.GetExtraction(ctx, req.ID, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ExtractionResponse{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		}
		return ExtractionResponse{}, err
	}
	return toResponse(*e), nil
}
