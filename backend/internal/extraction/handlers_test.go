package extraction_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/basselshurbaji/mr_bean/backend/internal/extraction"
	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type mockExtractionService struct {
	item *extraction.Extraction
	list []extraction.Extraction
	err  error
}

func (m *mockExtractionService) ListExtractions(_ context.Context, _ string, _, _ int) ([]extraction.Extraction, error) {
	return m.list, m.err
}
func (m *mockExtractionService) GetExtraction(_ context.Context, _, _ string) (*extraction.Extraction, error) {
	return m.item, m.err
}
func (m *mockExtractionService) CreateExtraction(_ context.Context, _ string, _ extraction.ExtractionParams) (*extraction.Extraction, error) {
	return m.item, m.err
}
func (m *mockExtractionService) UpdateExtraction(_ context.Context, _, _ string, _ extraction.ExtractionParams) (*extraction.Extraction, error) {
	return m.item, m.err
}
func (m *mockExtractionService) DeleteExtraction(_ context.Context, _, _ string) error { return m.err }

func authedCtx() context.Context {
	return principal.WithUserID(context.Background(), "user-1")
}

var sampleExtraction = makeExtraction("e1")

// --- ListExtractionsHandler ---

func TestListExtractionsHandler_Validate(t *testing.T) {
	h := extraction.NewListExtractionsHandler(nil)
	if err := h.Validate(extraction.ListExtractionsRequest{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListExtractionsHandler_Serve_EmptyList(t *testing.T) {
	svc := &mockExtractionService{list: []extraction.Extraction{}}
	h := extraction.NewListExtractionsHandler(svc)

	res, err := h.Serve(authedCtx(), extraction.ListExtractionsRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("expected empty list, got %d items", len(res))
	}
}

func TestListExtractionsHandler_Serve_ReturnsList(t *testing.T) {
	items := []extraction.Extraction{sampleExtraction, makeExtraction("e2")}
	svc := &mockExtractionService{list: items}
	h := extraction.NewListExtractionsHandler(svc)

	res, err := h.Serve(authedCtx(), extraction.ListExtractionsRequest{Limit: 10, Page: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 items, got %d", len(res))
	}
}

// --- GetExtractionHandler ---

func TestGetExtractionHandler_Validate(t *testing.T) {
	h := extraction.NewGetExtractionHandler(nil)

	if err := h.Validate(extraction.GetExtractionRequest{ID: "e1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Validate(extraction.GetExtractionRequest{}); err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestGetExtractionHandler_Serve_Success(t *testing.T) {
	svc := &mockExtractionService{item: &sampleExtraction}
	h := extraction.NewGetExtractionHandler(svc)

	res, err := h.Serve(authedCtx(), extraction.GetExtractionRequest{ID: "e1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != sampleExtraction.ID {
		t.Fatalf("expected id %q, got %q", sampleExtraction.ID, res.ID)
	}
}

func TestGetExtractionHandler_Serve_NotFound(t *testing.T) {
	svc := &mockExtractionService{err: extraction.ErrNotFound}
	h := extraction.NewGetExtractionHandler(svc)

	_, err := h.Serve(authedCtx(), extraction.GetExtractionRequest{ID: "missing"})
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", appErr.Code)
	}
}

// --- CreateExtractionHandler ---

func TestCreateExtractionHandler_Validate(t *testing.T) {
	h := extraction.NewCreateExtractionHandler(nil)

	valid := extraction.CreateExtractionRequest{
		BeanID: "bean-1", DoseIn: 18, YieldOut: 36, Time: 27, TargetTime: 27, GrindSize: 14,
	}
	cases := []struct {
		name    string
		req     extraction.CreateExtractionRequest
		wantErr bool
	}{
		{"valid", valid, false},
		{"missing bean_id", func() extraction.CreateExtractionRequest { r := valid; r.BeanID = ""; return r }(), true},
		{"dose_in zero", func() extraction.CreateExtractionRequest { r := valid; r.DoseIn = 0; return r }(), true},
		{"yield_out negative", func() extraction.CreateExtractionRequest { r := valid; r.YieldOut = -1; return r }(), true},
		{"time zero", func() extraction.CreateExtractionRequest { r := valid; r.Time = 0; return r }(), true},
		{"target_time zero", func() extraction.CreateExtractionRequest { r := valid; r.TargetTime = 0; return r }(), true},
		{"grind_size zero", func() extraction.CreateExtractionRequest { r := valid; r.GrindSize = 0; return r }(), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := h.Validate(tc.req)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCreateExtractionHandler_Serve_Success(t *testing.T) {
	svc := &mockExtractionService{item: &sampleExtraction}
	h := extraction.NewCreateExtractionHandler(svc)

	req := extraction.CreateExtractionRequest{
		BeanID: "bean-1", DoseIn: 18, YieldOut: 36, Time: 27, TargetTime: 27, GrindSize: 14,
	}
	res, err := h.Serve(authedCtx(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != sampleExtraction.ID {
		t.Fatalf("expected id %q, got %q", sampleExtraction.ID, res.ID)
	}
}

func TestCreateExtractionHandler_Serve_InvalidBean(t *testing.T) {
	svc := &mockExtractionService{err: extraction.ErrInvalidBean}
	h := extraction.NewCreateExtractionHandler(svc)

	req := extraction.CreateExtractionRequest{
		BeanID: "bad-bean", DoseIn: 18, YieldOut: 36, Time: 27, TargetTime: 27, GrindSize: 14,
	}
	_, err := h.Serve(authedCtx(), req)
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", appErr.Code)
	}
}

func TestCreateExtractionHandler_Serve_InvalidGear(t *testing.T) {
	svc := &mockExtractionService{err: extraction.ErrInvalidGear}
	h := extraction.NewCreateExtractionHandler(svc)

	req := extraction.CreateExtractionRequest{
		BeanID: "bean-1", DoseIn: 18, YieldOut: 36, Time: 27, TargetTime: 27, GrindSize: 14,
		GearIDs: []string{"bad-gear"},
	}
	_, err := h.Serve(authedCtx(), req)
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", appErr.Code)
	}
}

// --- UpdateExtractionHandler ---

func TestUpdateExtractionHandler_Validate(t *testing.T) {
	h := extraction.NewUpdateExtractionHandler(nil)

	valid := extraction.UpdateExtractionRequest{
		ID: "e1", BeanID: "bean-1", DoseIn: 18, YieldOut: 36, Time: 27, TargetTime: 27, GrindSize: 14,
	}
	if err := h.Validate(valid); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	missing := valid
	missing.BeanID = ""
	if err := h.Validate(missing); err == nil {
		t.Fatal("expected error for missing bean_id")
	}
}

func TestUpdateExtractionHandler_Serve_NotFound(t *testing.T) {
	svc := &mockExtractionService{err: extraction.ErrNotFound}
	h := extraction.NewUpdateExtractionHandler(svc)

	req := extraction.UpdateExtractionRequest{
		ID: "missing", BeanID: "bean-1", DoseIn: 18, YieldOut: 36, Time: 27, TargetTime: 27, GrindSize: 14,
	}
	_, err := h.Serve(authedCtx(), req)
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", appErr.Code)
	}
}

func TestUpdateExtractionHandler_Serve_InvalidBean(t *testing.T) {
	svc := &mockExtractionService{err: extraction.ErrInvalidBean}
	h := extraction.NewUpdateExtractionHandler(svc)

	req := extraction.UpdateExtractionRequest{
		ID: "e1", BeanID: "bad-bean", DoseIn: 18, YieldOut: 36, Time: 27, TargetTime: 27, GrindSize: 14,
	}
	_, err := h.Serve(authedCtx(), req)
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", appErr.Code)
	}
}

// --- DeleteExtractionHandler ---

func TestDeleteExtractionHandler_Validate(t *testing.T) {
	h := extraction.NewDeleteExtractionHandler(nil)

	if err := h.Validate(extraction.DeleteExtractionRequest{ID: "e1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Validate(extraction.DeleteExtractionRequest{}); err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestDeleteExtractionHandler_Serve_Success(t *testing.T) {
	svc := &mockExtractionService{}
	h := extraction.NewDeleteExtractionHandler(svc)

	if _, err := h.Serve(authedCtx(), extraction.DeleteExtractionRequest{ID: "e1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteExtractionHandler_Serve_NotFound(t *testing.T) {
	svc := &mockExtractionService{err: extraction.ErrNotFound}
	h := extraction.NewDeleteExtractionHandler(svc)

	_, err := h.Serve(authedCtx(), extraction.DeleteExtractionRequest{ID: "missing"})
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", appErr.Code)
	}
}
