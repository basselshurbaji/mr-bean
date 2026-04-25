package bean

import (
	"context"
	"errors"
	"net/http"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type DeleteBeanRequest struct {
	ID string `url:"id"`
}

type DeleteBeanHandler struct {
	svc BeanService
}

func NewDeleteBeanHandler(svc BeanService) *DeleteBeanHandler {
	return &DeleteBeanHandler{svc: svc}
}

func (h *DeleteBeanHandler) Method() string  { return "DELETE" }
func (h *DeleteBeanHandler) Pattern() string { return "/beans/{id}" }

func (h *DeleteBeanHandler) Middlewares() []middleware.Tag {
	return []middleware.Tag{middleware.TagAuthenticated}
}

func (h *DeleteBeanHandler) Validate(req DeleteBeanRequest) error {
	if req.ID == "" {
		return errors.New("id is required")
	}
	return nil
}

func (h *DeleteBeanHandler) Serve(ctx context.Context, req DeleteBeanRequest) (handler.NoContent, error) {
	userID, ok := principal.UserIDFromContext(ctx)
	if !ok {
		return handler.NoContent{}, errors.New("unauthorized")
	}
	if err := h.svc.DeleteBean(ctx, req.ID, userID); err != nil {
		if errors.Is(err, ErrNotFound) {
			return handler.NoContent{}, &handler.AppError{Code: http.StatusNotFound, Msg: "not found"}
		}
		return handler.NoContent{}, err
	}
	return handler.NoContent{}, nil
}
