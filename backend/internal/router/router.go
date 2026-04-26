package router

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/schema"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
)

var queryDecoder = schema.NewDecoder()
var urlDecoder = schema.NewDecoder()

func init() {
	urlDecoder.SetAliasTag("url")
}

type Router struct {
	chi chi.Router
}

func NewRouter() *Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	return &Router{chi: r}
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}

func Register[Req, Res any](r *Router, h handler.Handler[Req, Res]) {
	mws := middleware.Resolve(h.Middlewares())

	chain := r.chi
	if len(mws) > 0 {
		chain = r.chi.With(mws...)
	}

	chain.MethodFunc(h.Method(), h.Pattern(), func(w http.ResponseWriter, req *http.Request) {
		var body Req

		switch req.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse("invalid request body"))
				return
			}
		case http.MethodGet:
			if err := queryDecoder.Decode(&body, req.URL.Query()); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse("invalid query parameters"))
				return
			}
		case http.MethodDelete:
			// no body
		default:
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse("unsupported method"))
			return
		}

		// Inject chi URL path params into the request struct using `url` tags.
		if rctx := chi.RouteContext(req.Context()); rctx != nil {
			vals := make(url.Values)
			for i, k := range rctx.URLParams.Keys {
				vals.Set(k, rctx.URLParams.Values[i])
			}
			if len(vals) > 0 {
				if err := urlDecoder.Decode(&body, vals); err != nil {
					writeJSON(w, http.StatusBadRequest, errorResponse("invalid url parameters"))
					return
				}
			}
		}

		if err := h.Validate(body); err != nil {
			writeJSON(w, http.StatusUnprocessableEntity, errorResponse(err.Error()))
			return
		}

		res, err := h.Serve(req.Context(), body)
		if err != nil {
			if appErr, ok := err.(*handler.AppError); ok {
				writeJSON(w, appErr.Code, errorResponse(appErr.Msg))
				return
			}
			writeJSON(w, http.StatusInternalServerError, errorResponse(err.Error()))
			return
		}

		if _, ok := any(res).(handler.NoContent); ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		writeJSON(w, http.StatusOK, res)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func errorResponse(msg string) map[string]string {
	return map[string]string{"error": msg}
}
