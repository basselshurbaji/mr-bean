package extraction_test

import (
	"context"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/extraction"
)

type mockExtractionRepo struct {
	item *extraction.Extraction
	list []extraction.Extraction
	err  error
}

func (m *mockExtractionRepo) ListExtractions(_ context.Context, _ string, _, _ int) ([]extraction.Extraction, error) {
	return m.list, m.err
}
func (m *mockExtractionRepo) GetExtraction(_ context.Context, _, _ string) (*extraction.Extraction, error) {
	return m.item, m.err
}
func (m *mockExtractionRepo) CreateExtraction(_ context.Context, _ string, _ extraction.ExtractionParams) (*extraction.Extraction, error) {
	return m.item, m.err
}
func (m *mockExtractionRepo) UpdateExtraction(_ context.Context, _, _ string, _ extraction.ExtractionParams) (*extraction.Extraction, error) {
	return m.item, m.err
}
func (m *mockExtractionRepo) DeleteExtraction(_ context.Context, _, _ string) error { return m.err }

func makeExtraction(id string) extraction.Extraction {
	return extraction.Extraction{
		ID:     id,
		UserID: "user-1",
		BeanID: "bean-1",
		Bean:   extraction.BeanInfo{ID: "bean-1", Name: "Ethiopia"},
		DoseIn: 18, YieldOut: 36, Time: 27, TargetTime: 27, GrindSize: 14,
		Gear:      []extraction.GearInfo{},
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func TestExtractionService_ListExtractions_ReturnsList(t *testing.T) {
	items := []extraction.Extraction{makeExtraction("e1"), makeExtraction("e2")}
	svc := extraction.NewExtractionService(&mockExtractionRepo{list: items})

	got, err := svc.ListExtractions(context.Background(), "user-1", 20, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
}

func TestExtractionService_GetExtraction_ReturnsItem(t *testing.T) {
	e := makeExtraction("e1")
	svc := extraction.NewExtractionService(&mockExtractionRepo{item: &e})

	got, err := svc.GetExtraction(context.Background(), "e1", "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "e1" {
		t.Fatalf("expected id %q, got %q", "e1", got.ID)
	}
}

func TestExtractionService_CreateExtraction_ReturnsExtraction(t *testing.T) {
	e := makeExtraction("e1")
	svc := extraction.NewExtractionService(&mockExtractionRepo{item: &e})

	p := extraction.ExtractionParams{
		BeanID: "bean-1", DoseIn: 18, YieldOut: 36, Time: 27,
		TargetTime: 27, GrindSize: 14, GearIDs: []string{},
	}
	got, err := svc.CreateExtraction(context.Background(), "user-1", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "e1" {
		t.Fatalf("expected id %q, got %q", "e1", got.ID)
	}
}

func TestExtractionService_UpdateExtraction_NotFound(t *testing.T) {
	svc := extraction.NewExtractionService(&mockExtractionRepo{err: extraction.ErrNotFound})

	p := extraction.ExtractionParams{BeanID: "bean-1", DoseIn: 18, YieldOut: 36, Time: 27, TargetTime: 27, GrindSize: 14}
	_, err := svc.UpdateExtraction(context.Background(), "missing", "user-1", p)
	if err != extraction.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestExtractionService_DeleteExtraction_Success(t *testing.T) {
	svc := extraction.NewExtractionService(&mockExtractionRepo{})

	if err := svc.DeleteExtraction(context.Background(), "e1", "user-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExtractionService_DeleteExtraction_NotFound(t *testing.T) {
	svc := extraction.NewExtractionService(&mockExtractionRepo{err: extraction.ErrNotFound})

	err := svc.DeleteExtraction(context.Background(), "missing", "user-1")
	if err != extraction.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
