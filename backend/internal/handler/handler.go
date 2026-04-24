package handler

import (
	"context"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
)

// Handler defines a single HTTP endpoint. Req and Res are plain structs —
// no HTTP types leak into implementations.
type Handler[Req, Res any] interface {
	Method() string
	Pattern() string
	Middlewares() []middleware.Tag
	Validate(req Req) error
	Serve(ctx context.Context, req Req) (Res, error)
}

// AppError is returned by handlers to signal a specific HTTP status code
// instead of the default 500.
type AppError struct {
	Code int
	Msg  string
}

func (e *AppError) Error() string { return e.Msg }

// NoContent is the response type for handlers that return HTTP 204 with no body.
type NoContent struct{}
