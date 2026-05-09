package extraction

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type DeleteExtractionRequest struct {
	ID string `url:"id"`
}

type DeleteExtractionHandler struct {
	svc ExtractionService
}

func NewDeleteExtractionHandler(svc ExtractionService) *DeleteExtractionHandler {
	return &DeleteExtractionHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *DeleteExtractionHandler) Method() string { return "DELETE" }

// Pattern implements handler.Handler.
func (h *DeleteExtractionHandler) Pattern() string { return "/extractions/{id}" }

// Middlewares implements handler.Handler.
func (h *DeleteExtractionHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated, middleware.TagAppAuthenticated}
}

// Validate implements handler.Handler.
func (h *DeleteExtractionHandler) Validate(req DeleteExtractionRequest) error {
	if req.ID == "" {
		return errors.New("id is required")
	}
	return nil
}

// Serve implements handler.Handler.
func (h *DeleteExtractionHandler) Serve(ctx context.Context, req DeleteExtractionRequest) (handler.NoContent, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return handler.NoContent{}, errors.New("unauthorized")
	}
	if err := h.svc.DeleteExtraction(ctx, req.ID, userID); err != nil {
		if errors.Is(err, ErrNotFound) {
			return handler.NoContent{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		}
		return handler.NoContent{}, err
	}
	return handler.NoContent{}, nil
}
