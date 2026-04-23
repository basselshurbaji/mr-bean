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
