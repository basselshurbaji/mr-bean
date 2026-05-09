package extraction

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type UpdateExtractionRequest struct {
	ID          string   `url:"id"`
	BeanID      string   `json:"bean_id"`
	DoseIn      float64  `json:"dose_in"`
	YieldOut    float64  `json:"yield_out"`
	Time        float64  `json:"time"`
	TargetTime  float64  `json:"target_time"`
	GrindSize   float64  `json:"grind_size"`
	GearIDs     []string `json:"gear_ids"`
	PreInfusion bool     `json:"pre_infusion"`
	TastingNote *string  `json:"tasting_note"`
}

type UpdateExtractionHandler struct {
	svc ExtractionService
}

func NewUpdateExtractionHandler(svc ExtractionService) *UpdateExtractionHandler {
	return &UpdateExtractionHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *UpdateExtractionHandler) Method() string { return "PUT" }

// Pattern implements handler.Handler.
func (h *UpdateExtractionHandler) Pattern() string { return "/extractions/{id}" }

// Middlewares implements handler.Handler.
func (h *UpdateExtractionHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated, middleware.TagAppAuthenticated}
}

// Validate implements handler.Handler.
func (h *UpdateExtractionHandler) Validate(req UpdateExtractionRequest) error {
	return validateExtractionFields(req.BeanID, req.DoseIn, req.YieldOut, req.Time, req.TargetTime, req.GrindSize)
}

// Serve implements handler.Handler.
func (h *UpdateExtractionHandler) Serve(ctx context.Context, req UpdateExtractionRequest) (ExtractionResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return ExtractionResponse{}, errors.New("unauthorized")
	}
	gearIDs := req.GearIDs
	if gearIDs == nil {
		gearIDs = []string{}
	}
	e, err := h.svc.UpdateExtraction(ctx, req.ID, userID, ExtractionParams{
		BeanID:      req.BeanID,
		DoseIn:      req.DoseIn,
		YieldOut:    req.YieldOut,
		Time:        req.Time,
		TargetTime:  req.TargetTime,
		GrindSize:   req.GrindSize,
		PreInfusion: req.PreInfusion,
		TastingNote: req.TastingNote,
		GearIDs:     gearIDs,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return ExtractionResponse{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		case errors.Is(err, ErrInvalidBean), errors.Is(err, ErrInvalidGear):
			return ExtractionResponse{}, &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: err.Error()}
		}
		return ExtractionResponse{}, err
	}
	return toResponse(*e), nil
}
