# internal/handler

Contains only the `Handler[Req, Res any]` interface — the contract every endpoint implements.

```go
type Handler[Req, Res any] interface {
    Method() string
    Pattern() string
    Validate(req Req) error
    Serve(ctx context.Context, req Req) (Res, error)
}
```

This package has no dependencies other than `context`. Do not add HTTP types, middleware, or any feature code here.
