package gear_test

import (
	"context"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/gear"
)

// mockGearRepo stubs GearRepo for service tests.
type mockGearRepo struct {
	gearList []gear.GearItem
	gearItem *gear.GearItem
	station  *gear.Station
	stations []gear.Station
	err      error
}

func (m *mockGearRepo) ListGear(_ context.Context, _ string) ([]gear.GearItem, error) {
	return m.gearList, m.err
}
func (m *mockGearRepo) GetGear(_ context.Context, _, _ string) (*gear.GearItem, error) {
	return m.gearItem, m.err
}
func (m *mockGearRepo) CreateGear(_ context.Context, _ string, _ gear.CreateGearParams) (*gear.GearItem, error) {
	return m.gearItem, m.err
}
func (m *mockGearRepo) UpdateGear(_ context.Context, _, _ string, _ gear.UpdateGearParams) (*gear.GearItem, error) {
	return m.gearItem, m.err
}
func (m *mockGearRepo) DeleteGear(_ context.Context, _, _ string) error { return m.err }
func (m *mockGearRepo) ListStations(_ context.Context, _ string) ([]gear.Station, error) {
	return m.stations, m.err
}
func (m *mockGearRepo) CreateStation(_ context.Context, _, _ string, _ []string) (*gear.Station, error) {
	return m.station, m.err
}
func (m *mockGearRepo) UpdateStation(_ context.Context, _, _, _ string, _ []string) (*gear.Station, error) {
	return m.station, m.err
}
func (m *mockGearRepo) DeleteStation(_ context.Context, _, _ string) error { return m.err }

func makeGearItem(id string) gear.GearItem {
	return gear.GearItem{
		ID: id, UserID: "user-1", TypeID: "grinder", Name: "Test",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func TestGearService_CreateStation_EmptyGear(t *testing.T) {
	repo := &mockGearRepo{
		station: &gear.Station{ID: "s1", UserID: "user-1", Name: "Morning"},
	}
	svc := gear.NewGearService(repo)

	s, err := svc.CreateStation(context.Background(), "user-1", "Morning", []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Gear) != 0 {
		t.Fatalf("expected 0 gear, got %d", len(s.Gear))
	}
}

func TestGearService_CreateStation_UnownedGear(t *testing.T) {
	repo := &mockGearRepo{
		gearList: []gear.GearItem{makeGearItem("gear-1")},
	}
	svc := gear.NewGearService(repo)

	_, err := svc.CreateStation(context.Background(), "user-1", "Morning", []string{"gear-1", "bad-id"})
	if err == nil {
		t.Fatal("expected error for unowned gear")
	}
	if err != gear.ErrUnownedGear {
		t.Fatalf("expected ErrUnownedGear, got %v", err)
	}
}

func TestGearService_CreateStation_OwnedGear_OrderPreserved(t *testing.T) {
	g1 := makeGearItem("gear-1")
	g2 := makeGearItem("gear-2")
	repo := &mockGearRepo{
		gearList: []gear.GearItem{g1, g2},
		station:  &gear.Station{ID: "s1", UserID: "user-1", Name: "Morning"},
	}
	svc := gear.NewGearService(repo)

	s, err := svc.CreateStation(context.Background(), "user-1", "Morning", []string{"gear-2", "gear-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Gear) != 2 {
		t.Fatalf("expected 2 gear, got %d", len(s.Gear))
	}
	if s.Gear[0].ID != "gear-2" || s.Gear[1].ID != "gear-1" {
		t.Fatalf("order not preserved: %v %v", s.Gear[0].ID, s.Gear[1].ID)
	}
}

func TestGearService_UpdateStation_UnownedGear(t *testing.T) {
	repo := &mockGearRepo{
		gearList: []gear.GearItem{makeGearItem("gear-1")},
	}
	svc := gear.NewGearService(repo)

	_, err := svc.UpdateStation(context.Background(), "s1", "user-1", "Morning", []string{"not-owned"})
	if err != gear.ErrUnownedGear {
		t.Fatalf("expected ErrUnownedGear, got %v", err)
	}
}
