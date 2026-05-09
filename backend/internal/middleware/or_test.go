package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/basselshurbaji/mr_bean/backend/internal/middleware"
)

func rejectMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
}

func acceptMW(id string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Auth", id)
			next.ServeHTTP(w, r)
		})
	}
}

func okNext(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }

func TestOr_FirstAccepts(t *testing.T) {
	mw := middleware.Or(acceptMW("first"), rejectMW)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(http.HandlerFunc(okNext)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if rec.Header().Get("X-Auth") != "first" {
		t.Fatalf("want X-Auth=first, got %q", rec.Header().Get("X-Auth"))
	}
}

func TestOr_FirstRejectsSecondAccepts(t *testing.T) {
	mw := middleware.Or(rejectMW, acceptMW("second"))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(http.HandlerFunc(okNext)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if rec.Header().Get("X-Auth") != "second" {
		t.Fatalf("want X-Auth=second, got %q", rec.Header().Get("X-Auth"))
	}
}

func TestOr_BothReject(t *testing.T) {
	mw := middleware.Or(rejectMW, rejectMW)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(http.HandlerFunc(okNext)).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}
}
