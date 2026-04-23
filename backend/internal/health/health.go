package health

import (
	"context"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
)

type Request struct{}

type Response struct {
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Method() string              { return "GET" }
func (h *Handler) Pattern() string             { return "/health" }
func (h *Handler) Middlewares() []middleware.Tag { return nil }

func (h *Handler) Validate(_ Request) error { return nil }

func (h *Handler) Serve(_ context.Context, _ Request) (Response, error) {
	return Response{
		Timestamp: time.Now().Unix(),
		Status:    "healthy",
	}, nil
}