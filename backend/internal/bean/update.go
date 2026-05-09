package bean

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type UpdateBeanRequest struct {
	ID           string  `url:"id"`
	Name         string  `json:"name"`
	Roaster      *string `json:"roaster"`
	Origin       *string `json:"origin"`
	Process      *string `json:"process"`
	RoastLevel   *string `json:"roast_level"`
	TastingNotes *string `json:"tasting_notes"`
	Notes        *string `json:"notes"`
}

type UpdateBeanHandler struct {
	svc BeanService
}

func NewUpdateBeanHandler(svc BeanService) *UpdateBeanHandler {
	return &UpdateBeanHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *UpdateBeanHandler) Method() string  { return "PUT" }
// Pattern implements handler.Handler.
func (h *UpdateBeanHandler) Pattern() string { return "/beans/{id}" }

// Middlewares implements handler.Handler.
func (h *UpdateBeanHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated, middleware.TagAppAuthenticated}
}

// Validate implements handler.Handler.
func (h *UpdateBeanHandler) Validate(req UpdateBeanRequest) error {
	return validateBeanFields(req.Name, req.Process, req.RoastLevel)
}

// Serve implements handler.Handler.
func (h *UpdateBeanHandler) Serve(ctx context.Context, req UpdateBeanRequest) (BeanResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return BeanResponse{}, errors.New("unauthorized")
	}
	b, err := h.svc.UpdateBean(ctx, req.ID, userID, BeanParams{
		Name:         req.Name,
		Roaster:      req.Roaster,
		Origin:       req.Origin,
		Process:      req.Process,
		RoastLevel:   req.RoastLevel,
		TastingNotes: req.TastingNotes,
		Notes:        req.Notes,
	})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return BeanResponse{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		}
		return BeanResponse{}, err
	}
	return beanToResponse(*b), nil
}
