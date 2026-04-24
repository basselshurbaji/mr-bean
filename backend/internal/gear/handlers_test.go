package gear_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/gear"
	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

// mockGearService stubs GearService for handler tests.
type mockGearService struct {
	gearItem *gear.GearItem
	gearList []gear.GearItem
	station  *gear.Station
	stations []gear.Station
	err      error
}

func (m *mockGearService) ListGear(_ context.Context, _ string) ([]gear.GearItem, error) {
	return m.gearList, m.err
}
func (m *mockGearService) GetGear(_ context.Context, _, _ string) (*gear.GearItem, error) {
	return m.gearItem, m.err
}
func (m *mockGearService) CreateGear(_ context.Context, _ string, _ gear.CreateGearParams) (*gear.GearItem, error) {
	return m.gearItem, m.err
}
func (m *mockGearService) UpdateGear(_ context.Context, _, _ string, _ gear.UpdateGearParams) (*gear.GearItem, error) {
	return m.gearItem, m.err
}
func (m *mockGearService) DeleteGear(_ context.Context, _, _ string) error { return m.err }
func (m *mockGearService) ListStations(_ context.Context, _ string) ([]gear.Station, error) {
	return m.stations, m.err
}
func (m *mockGearService) CreateStation(_ context.Context, _, _ string, _ []string) (*gear.Station, error) {
	return m.station, m.err
}
func (m *mockGearService) UpdateStation(_ context.Context, _, _, _ string, _ []string) (*gear.Station, error) {
	return m.station, m.err
}
func (m *mockGearService) DeleteStation(_ context.Context, _, _ string) error { return m.err }

func authedCtx() context.Context {
	return principal.WithUserID(context.Background(), "user-1")
}

var sampleGear = gear.GearItem{
	ID:        "gear-1",
	UserID:    "user-1",
	TypeID:    "grinder",
	Name:      "Niche Zero",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

// --- CreateGearHandler ---

func TestCreateGearHandler_Validate(t *testing.T) {
	h := gear.NewCreateGearHandler(nil)

	cases := []struct {
		name    string
		req     gear.CreateGearRequest
		wantErr bool
	}{
		{"valid", gear.CreateGearRequest{TypeID: "grinder", Name: "Niche"}, false},
		{"missing type_id", gear.CreateGearRequest{Name: "Niche"}, true},
		{"invalid type_id", gear.CreateGearRequest{TypeID: "unknown", Name: "Niche"}, true},
		{"missing name", gear.CreateGearRequest{TypeID: "grinder"}, true},
		{"bad year", gear.CreateGearRequest{TypeID: "grinder", Name: "X", Year: strPtr("99")}, true},
		{"good year", gear.CreateGearRequest{TypeID: "grinder", Name: "X", Year: strPtr("2022")}, false},
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

func TestCreateGearHandler_Serve_Success(t *testing.T) {
	svc := &mockGearService{gearItem: &sampleGear}
	h := gear.NewCreateGearHandler(svc)

	res, err := h.Serve(authedCtx(), gear.CreateGearRequest{TypeID: "grinder", Name: "Niche"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != sampleGear.ID {
		t.Fatalf("expected id %q, got %q", sampleGear.ID, res.ID)
	}
}

// --- UpdateGearHandler ---

func TestUpdateGearHandler_Validate(t *testing.T) {
	h := gear.NewUpdateGearHandler(nil)

	if err := h.Validate(gear.UpdateGearRequest{TypeID: "scale", Name: "Decent"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Validate(gear.UpdateGearRequest{TypeID: "scale"}); err == nil {
		t.Fatal("expected error for missing name")
	}
	if err := h.Validate(gear.UpdateGearRequest{Name: "X"}); err == nil {
		t.Fatal("expected error for missing type_id")
	}
}

// --- GetGearHandler ---

func TestGetGearHandler_Serve_NotFound(t *testing.T) {
	svc := &mockGearService{err: gear.ErrNotFound}
	h := gear.NewGetGearHandler(svc)

	_, err := h.Serve(authedCtx(), gear.GetGearRequest{ID: "missing"})
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", appErr.Code)
	}
}

// --- DeleteGearHandler ---

func TestDeleteGearHandler_Validate(t *testing.T) {
	h := gear.NewDeleteGearHandler(nil)

	if err := h.Validate(gear.DeleteGearRequest{ID: "x"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Validate(gear.DeleteGearRequest{}); err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestDeleteGearHandler_Serve_NotFound(t *testing.T) {
	svc := &mockGearService{err: gear.ErrNotFound}
	h := gear.NewDeleteGearHandler(svc)

	_, err := h.Serve(authedCtx(), gear.DeleteGearRequest{ID: "missing"})
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", appErr.Code)
	}
}

// --- CreateStationHandler ---

func TestCreateStationHandler_Validate(t *testing.T) {
	h := gear.NewCreateStationHandler(nil)

	if err := h.Validate(gear.CreateStationRequest{Name: "Morning"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Validate(gear.CreateStationRequest{}); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestCreateStationHandler_Serve_UnownedGear(t *testing.T) {
	svc := &mockGearService{err: gear.ErrUnownedGear}
	h := gear.NewCreateStationHandler(svc)

	_, err := h.Serve(authedCtx(), gear.CreateStationRequest{Name: "Morning", GearIDs: []string{"bad-id"}})
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", appErr.Code)
	}
}

// --- UpdateStationHandler ---

func TestUpdateStationHandler_Serve_NotFound(t *testing.T) {
	svc := &mockGearService{err: gear.ErrNotFound}
	h := gear.NewUpdateStationHandler(svc)

	_, err := h.Serve(authedCtx(), gear.UpdateStationRequest{ID: "x", Name: "Y"})
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", appErr.Code)
	}
}

// --- DeleteStationHandler ---

func TestDeleteStationHandler_Validate(t *testing.T) {
	h := gear.NewDeleteStationHandler(nil)

	if err := h.Validate(gear.DeleteStationRequest{ID: "x"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Validate(gear.DeleteStationRequest{}); err == nil {
		t.Fatal("expected error for missing id")
	}
}

func strPtr(s string) *string { return &s }
