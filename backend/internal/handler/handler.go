package handler

// Handler defines a single HTTP endpoint. Req and Res are plain structs —
// no HTTP types leak into implementations.
type Handler[Req, Res any] interface {
	Method() string
	Pattern() string
	Validate(req Req) error
	Serve(req Req) (Res, error)
}
