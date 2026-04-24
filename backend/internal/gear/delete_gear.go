package gear

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type DeleteGearRequest struct {
	ID string `url:"id"`
}

type DeleteGearHandler struct {
	svc GearService
}

func NewDeleteGearHandler(svc GearService) *DeleteGearHandler {
	return &DeleteGearHandler{svc: svc}
}

func (h *DeleteGearHandler) Method() string  { return "DELETE" }
func (h *DeleteGearHandler) Pattern() string { return "/gear/{id}" }

func (h *DeleteGearHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

func (h *DeleteGearHandler) Validate(req DeleteGearRequest) error {
	if req.ID == "" {
		return errors.New("id is required")
	}
	return nil
}

func (h *DeleteGearHandler) Serve(ctx context.Context, req DeleteGearRequest) (handler.NoContent, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return handler.NoContent{}, errors.New("unauthorized")
	}
	if err := h.svc.DeleteGear(ctx, req.ID, userID); err != nil {
		if errors.Is(err, ErrNotFound) {
			return handler.NoContent{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		}
		return handler.NoContent{}, err
	}
	return handler.NoContent{}, nil
}
