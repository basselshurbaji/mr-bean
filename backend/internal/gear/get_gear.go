package gear

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type GetGearRequest struct {
	ID string `url:"id"`
}

type GetGearHandler struct {
	svc GearService
}

func NewGetGearHandler(svc GearService) *GetGearHandler {
	return &GetGearHandler{svc: svc}
}

func (h *GetGearHandler) Method() string  { return "GET" }
func (h *GetGearHandler) Pattern() string { return "/gear/{id}" }

func (h *GetGearHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

func (h *GetGearHandler) Validate(req GetGearRequest) error {
	if req.ID == "" {
		return errors.New("id is required")
	}
	return nil
}

func (h *GetGearHandler) Serve(ctx context.Context, req GetGearRequest) (GearResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return GearResponse{}, errors.New("unauthorized")
	}
	g, err := h.svc.GetGear(ctx, req.ID, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return GearResponse{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		}
		return GearResponse{}, err
	}
	return gearToResponse(*g), nil
}
