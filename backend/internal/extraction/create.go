package extraction

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type CreateExtractionRequest struct {
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

type CreateExtractionHandler struct {
	svc ExtractionService
}

func NewCreateExtractionHandler(svc ExtractionService) *CreateExtractionHandler {
	return &CreateExtractionHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *CreateExtractionHandler) Method() string { return "POST" }

// Pattern implements handler.Handler.
func (h *CreateExtractionHandler) Pattern() string { return "/extractions" }

// Middlewares implements handler.Handler.
func (h *CreateExtractionHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated, middleware.TagAppAuthenticated}
}

// Validate implements handler.Handler.
func (h *CreateExtractionHandler) Validate(req CreateExtractionRequest) error {
	return validateExtractionFields(req.BeanID, req.DoseIn, req.YieldOut, req.Time, req.TargetTime, req.GrindSize)
}

// Serve implements handler.Handler.
func (h *CreateExtractionHandler) Serve(ctx context.Context, req CreateExtractionRequest) (ExtractionResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return ExtractionResponse{}, errors.New("unauthorized")
	}
	gearIDs := req.GearIDs
	if gearIDs == nil {
		gearIDs = []string{}
	}
	e, err := h.svc.CreateExtraction(ctx, userID, ExtractionParams{
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
		if errors.Is(err, ErrInvalidBean) || errors.Is(err, ErrInvalidGear) {
			return ExtractionResponse{}, &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: err.Error()}
		}
		return ExtractionResponse{}, err
	}
	return toResponse(*e), nil
}

func validateExtractionFields(beanID string, doseIn, yieldOut, extractionTime, targetTime, grindSize float64) error {
	if beanID == "" {
		return &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: "bean_id is required"}
	}
	if doseIn <= 0 {
		return &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: "dose_in must be greater than 0"}
	}
	if yieldOut <= 0 {
		return &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: "yield_out must be greater than 0"}
	}
	if extractionTime <= 0 {
		return &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: "time must be greater than 0"}
	}
	if targetTime <= 0 {
		return &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: "target_time must be greater than 0"}
	}
	if grindSize <= 0 {
		return &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: "grind_size must be greater than 0"}
	}
	return nil
}
