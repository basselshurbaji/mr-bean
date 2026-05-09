package gear

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type UpdateGearRequest struct {
	ID     string  `url:"id"`
	TypeID string  `json:"type_id"`
	Name   string  `json:"name"`
	Brand  *string `json:"brand"`
	Model  *string `json:"model"`
	Year   *string `json:"year"`
	Notes  *string `json:"notes"`
}

type UpdateGearHandler struct {
	svc GearService
}

func NewUpdateGearHandler(svc GearService) *UpdateGearHandler {
	return &UpdateGearHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *UpdateGearHandler) Method() string  { return "PUT" }
// Pattern implements handler.Handler.
func (h *UpdateGearHandler) Pattern() string { return "/gear/{id}" }

// Middlewares implements handler.Handler.
func (h *UpdateGearHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAnyAuthenticated}
}

// Validate implements handler.Handler.
func (h *UpdateGearHandler) Validate(req UpdateGearRequest) error {
	return validateGearFields(req.TypeID, req.Name, req.Year)
}

// Serve implements handler.Handler.
func (h *UpdateGearHandler) Serve(ctx context.Context, req UpdateGearRequest) (GearResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return GearResponse{}, errors.New("unauthorized")
	}
	g, err := h.svc.UpdateGear(ctx, req.ID, userID, UpdateGearParams{
		TypeID: req.TypeID,
		Name:   req.Name,
		Brand:  req.Brand,
		Model:  req.Model,
		Year:   req.Year,
		Notes:  req.Notes,
	})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return GearResponse{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		}
		return GearResponse{}, err
	}
	return gearToResponse(*g), nil
}
