package handler_test

import (
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
)

func TestHealthHandler(t *testing.T) {
	h := handler.NewHealthHandler()

	if h.Method() != "GET" {
		t.Errorf("expected method GET, got %s", h.Method())
	}

	if h.Pattern() != "/health" {
		t.Errorf("expected pattern /health, got %s", h.Pattern())
	}

	if err := h.Validate(handler.HealthRequest{}); err != nil {
		t.Errorf("expected nil validation error, got %v", err)
	}

	before := time.Now().Unix()
	res, err := h.Serve(handler.HealthRequest{})
	after := time.Now().Unix()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "healthy" {
		t.Errorf("expected status healthy, got %s", res.Status)
	}
	if res.Timestamp < before || res.Timestamp > after {
		t.Errorf("timestamp %d out of expected range [%d, %d]", res.Timestamp, before, after)
	}
}
