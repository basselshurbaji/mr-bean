package gear

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type UpdateStationRequest struct {
	ID      string   `url:"id"`
	Name    string   `json:"name"`
	GearIDs []string `json:"gear_ids"`
}

type UpdateStationHandler struct {
	svc GearService
}

func NewUpdateStationHandler(svc GearService) *UpdateStationHandler {
	return &UpdateStationHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *UpdateStationHandler) Method() string  { return "PUT" }
// Pattern implements handler.Handler.
func (h *UpdateStationHandler) Pattern() string { return "/stations/{id}" }

// Middlewares implements handler.Handler.
func (h *UpdateStationHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

// Validate implements handler.Handler.
func (h *UpdateStationHandler) Validate(req UpdateStationRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// Serve implements handler.Handler.
func (h *UpdateStationHandler) Serve(ctx context.Context, req UpdateStationRequest) (StationResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return StationResponse{}, errors.New("unauthorized")
	}
	gearIDs := req.GearIDs
	if gearIDs == nil {
		gearIDs = []string{}
	}
	s, err := h.svc.UpdateStation(ctx, req.ID, userID, req.Name, gearIDs)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return StationResponse{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		}
		if errors.Is(err, ErrUnownedGear) {
			return StationResponse{}, &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: err.Error()}
		}
		return StationResponse{}, err
	}
	return stationToResponse(*s), nil
}
