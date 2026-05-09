package gear

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type CreateStationRequest struct {
	Name    string   `json:"name"`
	GearIDs []string `json:"gear_ids"`
}

type CreateStationHandler struct {
	svc GearService
}

func NewCreateStationHandler(svc GearService) *CreateStationHandler {
	return &CreateStationHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *CreateStationHandler) Method() string  { return "POST" }
// Pattern implements handler.Handler.
func (h *CreateStationHandler) Pattern() string { return "/stations" }

// Middlewares implements handler.Handler.
func (h *CreateStationHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAnyAuthenticated}
}

// Validate implements handler.Handler.
func (h *CreateStationHandler) Validate(req CreateStationRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// Serve implements handler.Handler.
func (h *CreateStationHandler) Serve(ctx context.Context, req CreateStationRequest) (StationResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return StationResponse{}, errors.New("unauthorized")
	}
	gearIDs := req.GearIDs
	if gearIDs == nil {
		gearIDs = []string{}
	}
	s, err := h.svc.CreateStation(ctx, userID, req.Name, gearIDs)
	if err != nil {
		if errors.Is(err, ErrUnownedGear) {
			return StationResponse{}, &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: err.Error()}
		}
		return StationResponse{}, err
	}
	return stationToResponse(*s), nil
}
