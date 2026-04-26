package gear

import (
	"context"
	"errors"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type ListStationsRequest struct{}

type ListStationsHandler struct {
	svc GearService
}

func NewListStationsHandler(svc GearService) *ListStationsHandler {
	return &ListStationsHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *ListStationsHandler) Method() string  { return "GET" }
// Pattern implements handler.Handler.
func (h *ListStationsHandler) Pattern() string { return "/stations" }

// Middlewares implements handler.Handler.
func (h *ListStationsHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

// Validate implements handler.Handler.
func (h *ListStationsHandler) Validate(_ ListStationsRequest) error { return nil }

// Serve implements handler.Handler.
func (h *ListStationsHandler) Serve(ctx context.Context, _ ListStationsRequest) ([]StationResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	stations, err := h.svc.ListStations(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make([]StationResponse, len(stations))
	for i, s := range stations {
		res[i] = stationToResponse(s)
	}
	return res, nil
}
