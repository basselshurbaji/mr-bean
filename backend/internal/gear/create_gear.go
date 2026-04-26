package gear

import (
	"context"
	"errors"
	"strconv"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

var validTypeIDs = map[string]bool{
	"machine":     true,
	"grinder":     true,
	"scale":       true,
	"portafilter": true,
	"tamper":      true,
	"distributor": true,
	"wdt":         true,
	"basket":      true,
	"puckscreen":  true,
	"other":       true,
}

type CreateGearRequest struct {
	TypeID string  `json:"type_id"`
	Name   string  `json:"name"`
	Brand  *string `json:"brand"`
	Model  *string `json:"model"`
	Year   *string `json:"year"`
	Notes  *string `json:"notes"`
}

type CreateGearHandler struct {
	svc GearService
}

func NewCreateGearHandler(svc GearService) *CreateGearHandler {
	return &CreateGearHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *CreateGearHandler) Method() string  { return "POST" }
// Pattern implements handler.Handler.
func (h *CreateGearHandler) Pattern() string { return "/gear" }

// Middlewares implements handler.Handler.
func (h *CreateGearHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

// Validate implements handler.Handler.
func (h *CreateGearHandler) Validate(req CreateGearRequest) error {
	return validateGearFields(req.TypeID, req.Name, req.Year)
}

// Serve implements handler.Handler.
func (h *CreateGearHandler) Serve(ctx context.Context, req CreateGearRequest) (GearResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return GearResponse{}, errors.New("unauthorized")
	}
	g, err := h.svc.CreateGear(ctx, userID, CreateGearParams(req))
	if err != nil {
		return GearResponse{}, err
	}
	return gearToResponse(*g), nil
}

func validateGearFields(typeID, name string, year *string) error {
	if typeID == "" {
		return errors.New("type_id is required")
	}
	if !validTypeIDs[typeID] {
		return errors.New("invalid type_id")
	}
	if name == "" {
		return errors.New("name is required")
	}
	if year != nil {
		if _, err := strconv.Atoi(*year); err != nil || len(*year) != 4 {
			return errors.New("year must be a 4-digit string")
		}
	}
	return nil
}
