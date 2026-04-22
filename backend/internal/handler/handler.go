package handler

import "context"

// Handler defines a single HTTP endpoint. Req and Res are plain structs —
// no HTTP types leak into implementations.
type Handler[Req, Res any] interface {
	Method() string
	Pattern() string
	Validate(req Req) error
	Serve(ctx context.Context, req Req) (Res, error)
}
