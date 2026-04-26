package gear

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type DeleteStationRequest struct {
	ID string `url:"id"`
}

type DeleteStationHandler struct {
	svc GearService
}

func NewDeleteStationHandler(svc GearService) *DeleteStationHandler {
	return &DeleteStationHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *DeleteStationHandler) Method() string  { return "DELETE" }
// Pattern implements handler.Handler.
func (h *DeleteStationHandler) Pattern() string { return "/stations/{id}" }

// Middlewares implements handler.Handler.
func (h *DeleteStationHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

// Validate implements handler.Handler.
func (h *DeleteStationHandler) Validate(req DeleteStationRequest) error {
	if req.ID == "" {
		return errors.New("id is required")
	}
	return nil
}

// Serve implements handler.Handler.
func (h *DeleteStationHandler) Serve(ctx context.Context, req DeleteStationRequest) (handler.NoContent, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return handler.NoContent{}, errors.New("unauthorized")
	}
	if err := h.svc.DeleteStation(ctx, req.ID, userID); err != nil {
		if errors.Is(err, ErrNotFound) {
			return handler.NoContent{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		}
		return handler.NoContent{}, err
	}
	return handler.NoContent{}, nil
}
