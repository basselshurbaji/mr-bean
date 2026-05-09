package bean

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

var validProcesses = map[string]bool{
	"washed":    true,
	"natural":   true,
	"honey":     true,
	"anaerobic": true,
	"other":     true,
}

var validRoastLevels = map[string]bool{
	"light":       true,
	"medium_light": true,
	"medium":      true,
	"medium_dark": true,
	"dark":        true,
}

type CreateBeanRequest struct {
	Name         string  `json:"name"`
	Roaster      *string `json:"roaster"`
	Origin       *string `json:"origin"`
	Process      *string `json:"process"`
	RoastLevel   *string `json:"roast_level"`
	TastingNotes *string `json:"tasting_notes"`
	Notes        *string `json:"notes"`
}

type CreateBeanHandler struct {
	svc BeanService
}

func NewCreateBeanHandler(svc BeanService) *CreateBeanHandler {
	return &CreateBeanHandler{svc: svc}
}

// Method implements handler.Handler.
func (h *CreateBeanHandler) Method() string  { return "POST" }
// Pattern implements handler.Handler.
func (h *CreateBeanHandler) Pattern() string { return "/beans" }

// Middlewares implements handler.Handler.
func (h *CreateBeanHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated, middleware.TagAppAuthenticated}
}

// Validate implements handler.Handler.
func (h *CreateBeanHandler) Validate(req CreateBeanRequest) error {
	return validateBeanFields(req.Name, req.Process, req.RoastLevel)
}

// Serve implements handler.Handler.
func (h *CreateBeanHandler) Serve(ctx context.Context, req CreateBeanRequest) (BeanResponse, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return BeanResponse{}, errors.New("unauthorized")
	}
	b, err := h.svc.CreateBean(ctx, userID, BeanParams(req))
	if err != nil {
		return BeanResponse{}, err
	}
	return beanToResponse(*b), nil
}

func validateBeanFields(name string, process, roastLevel *string) error {
	if name == "" {
		return &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: "name is required"}
	}
	if process != nil && !validProcesses[*process] {
		return &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: "invalid process"}
	}
	if roastLevel != nil && !validRoastLevels[*roastLevel] {
		return &handler.AppError{Code: http.StatusUnprocessableEntity, Msg: "invalid roast_level"}
	}
	return nil
}
