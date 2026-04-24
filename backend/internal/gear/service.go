package gear

import (
	"context"
)

type GearService interface {
	ListGear(ctx context.Context, userID string) ([]GearItem, error)
	GetGear(ctx context.Context, id, userID string) (*GearItem, error)
	CreateGear(ctx context.Context, userID string, p CreateGearParams) (*GearItem, error)
	UpdateGear(ctx context.Context, id, userID string, p UpdateGearParams) (*GearItem, error)
	DeleteGear(ctx context.Context, id, userID string) error

	ListStations(ctx context.Context, userID string) ([]Station, error)
	CreateStation(ctx context.Context, userID, name string, gearIDs []string) (*Station, error)
	UpdateStation(ctx context.Context, id, userID, name string, gearIDs []string) (*Station, error)
	DeleteStation(ctx context.Context, id, userID string) error
}

type gearService struct {
	repo GearRepo
}

func NewGearService(repo GearRepo) GearService {
	return &gearService{repo: repo}
}

func (s *gearService) ListGear(ctx context.Context, userID string) ([]GearItem, error) {
	return s.repo.ListGear(ctx, userID)
}

func (s *gearService) GetGear(ctx context.Context, id, userID string) (*GearItem, error) {
	return s.repo.GetGear(ctx, id, userID)
}

func (s *gearService) CreateGear(ctx context.Context, userID string, p CreateGearParams) (*GearItem, error) {
	return s.repo.CreateGear(ctx, userID, p)
}

func (s *gearService) UpdateGear(ctx context.Context, id, userID string, p UpdateGearParams) (*GearItem, error) {
	return s.repo.UpdateGear(ctx, id, userID, p)
}

func (s *gearService) DeleteGear(ctx context.Context, id, userID string) error {
	return s.repo.DeleteGear(ctx, id, userID)
}

func (s *gearService) ListStations(ctx context.Context, userID string) ([]Station, error) {
	return s.repo.ListStations(ctx, userID)
}

func (s *gearService) CreateStation(ctx context.Context, userID, name string, gearIDs []string) (*Station, error) {
	if err := s.validateGearOwnership(ctx, userID, gearIDs); err != nil {
		return nil, err
	}
	station, err := s.repo.CreateStation(ctx, userID, name, gearIDs)
	if err != nil {
		return nil, err
	}
	station.Gear, err = s.orderedGear(ctx, userID, gearIDs)
	if err != nil {
		return nil, err
	}
	return station, nil
}

func (s *gearService) UpdateStation(ctx context.Context, id, userID, name string, gearIDs []string) (*Station, error) {
	if err := s.validateGearOwnership(ctx, userID, gearIDs); err != nil {
		return nil, err
	}
	station, err := s.repo.UpdateStation(ctx, id, userID, name, gearIDs)
	if err != nil {
		return nil, err
	}
	station.Gear, err = s.orderedGear(ctx, userID, gearIDs)
	if err != nil {
		return nil, err
	}
	return station, nil
}

func (s *gearService) DeleteStation(ctx context.Context, id, userID string) error {
	return s.repo.DeleteStation(ctx, id, userID)
}

// validateGearOwnership confirms every id in gearIDs exists and belongs to userID.
func (s *gearService) validateGearOwnership(ctx context.Context, userID string, gearIDs []string) error {
	if len(gearIDs) == 0 {
		return nil
	}
	allGear, err := s.repo.ListGear(ctx, userID)
	if err != nil {
		return err
	}
	owned := make(map[string]bool, len(allGear))
	for _, g := range allGear {
		owned[g.ID] = true
	}
	for _, id := range gearIDs {
		if !owned[id] {
			return ErrUnownedGear
		}
	}
	return nil
}

// orderedGear returns gear items in the order specified by gearIDs.
func (s *gearService) orderedGear(ctx context.Context, userID string, gearIDs []string) ([]GearItem, error) {
	if len(gearIDs) == 0 {
		return []GearItem{}, nil
	}
	allGear, err := s.repo.ListGear(ctx, userID)
	if err != nil {
		return nil, err
	}
	index := make(map[string]GearItem, len(allGear))
	for _, g := range allGear {
		index[g.ID] = g
	}
	result := make([]GearItem, len(gearIDs))
	for i, id := range gearIDs {
		result[i] = index[id]
	}
	return result, nil
}
