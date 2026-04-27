package extraction

import "context"

type ExtractionService interface {
	ListExtractions(ctx context.Context, userID string, limit, page int) ([]Extraction, error)
	GetExtraction(ctx context.Context, id, userID string) (*Extraction, error)
	CreateExtraction(ctx context.Context, userID string, p ExtractionParams) (*Extraction, error)
	UpdateExtraction(ctx context.Context, id, userID string, p ExtractionParams) (*Extraction, error)
	DeleteExtraction(ctx context.Context, id, userID string) error
}

type extractionService struct {
	repo ExtractionRepo
}

func NewExtractionService(repo ExtractionRepo) ExtractionService {
	return &extractionService{repo: repo}
}

// ListExtractions implements ExtractionService.
func (s *extractionService) ListExtractions(ctx context.Context, userID string, limit, page int) ([]Extraction, error) {
	return s.repo.ListExtractions(ctx, userID, limit, page)
}

// GetExtraction implements ExtractionService.
func (s *extractionService) GetExtraction(ctx context.Context, id, userID string) (*Extraction, error) {
	return s.repo.GetExtraction(ctx, id, userID)
}

// CreateExtraction implements ExtractionService.
func (s *extractionService) CreateExtraction(ctx context.Context, userID string, p ExtractionParams) (*Extraction, error) {
	return s.repo.CreateExtraction(ctx, userID, p)
}

// UpdateExtraction implements ExtractionService.
func (s *extractionService) UpdateExtraction(ctx context.Context, id, userID string, p ExtractionParams) (*Extraction, error) {
	return s.repo.UpdateExtraction(ctx, id, userID, p)
}

// DeleteExtraction implements ExtractionService.
func (s *extractionService) DeleteExtraction(ctx context.Context, id, userID string) error {
	return s.repo.DeleteExtraction(ctx, id, userID)
}
