package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/schema"
)

var queryDecoder = schema.NewDecoder()

// Route is an opaque value produced by Adapt. It carries the full HTTP lifecycle
// for one handler without exposing any router or HTTP types.
type Route struct {
	register func(chi.Router)
}

// Adapt wraps a Handler into a Route. chi and net/http are fully contained here.
func Adapt[Req, Res any](h Handler[Req, Res]) Route {
	return Route{
		register: func(r chi.Router) {
			r.MethodFunc(h.Method(), h.Pattern(), func(w http.ResponseWriter, req *http.Request) {
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
				default:
					writeJSON(w, http.StatusMethodNotAllowed, errorResponse("unsupported method"))
					return
				}

				if err := h.Validate(body); err != nil {
					writeJSON(w, http.StatusUnprocessableEntity, errorResponse(err.Error()))
					return
				}

				res, err := h.Serve(body)
				if err != nil {
					writeJSON(w, http.StatusInternalServerError, errorResponse(err.Error()))
					return
				}

				writeJSON(w, http.StatusOK, res)
			})
		},
	}
}

// Register mounts a single Route onto the router.
func Register(r chi.Router, route Route) {
	route.register(r)
}

// NewRouter creates the base chi router with standard middleware.
func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func errorResponse(msg string) map[string]string {
	return map[string]string{"error": msg}
}
