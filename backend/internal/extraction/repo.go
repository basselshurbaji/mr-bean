package extraction

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "github.com/basselshurbaji/mr_bean/backend/db/sqlc"
)

var ErrNotFound   = errors.New("not found")
var ErrInvalidBean = errors.New("bean_id is unknown or inaccessible")
var ErrInvalidGear = errors.New("one or more gear_ids are unknown or unowned")

type BeanInfo struct {
	ID      string
	Name    string
	Roaster *string
	Roast   *string
}

type GearInfo struct {
	ID     string
	TypeID string
	Name   string
}

type Extraction struct {
	ID          string
	UserID      string
	BeanID      string
	Bean        BeanInfo
	DoseIn      float64
	YieldOut    float64
	Time        float64
	TargetTime  float64
	GrindSize   float64
	PreInfusion bool
	TastingNote *string
	Gear        []GearInfo
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ExtractionParams struct {
	BeanID      string
	DoseIn      float64
	YieldOut    float64
	Time        float64
	TargetTime  float64
	GrindSize   float64
	PreInfusion bool
	TastingNote *string
	GearIDs     []string
}

type ExtractionRepo interface {
	ListExtractions(ctx context.Context, userID string, limit, page int) ([]Extraction, error)
	GetExtraction(ctx context.Context, id, userID string) (*Extraction, error)
	CreateExtraction(ctx context.Context, userID string, p ExtractionParams) (*Extraction, error)
	UpdateExtraction(ctx context.Context, id, userID string, p ExtractionParams) (*Extraction, error)
	DeleteExtraction(ctx context.Context, id, userID string) error
}

type pgExtractionRepo struct {
	db *sql.DB
	q  *db.Queries
}

func NewPgExtractionRepo(d *sql.DB) ExtractionRepo {
	return &pgExtractionRepo{db: d, q: db.New(d)}
}

// ListExtractions implements ExtractionRepo.
func (r *pgExtractionRepo) ListExtractions(ctx context.Context, userID string, limit, page int) ([]Extraction, error) {
	offset := (page - 1) * limit
	rows, err := r.q.ListExtractionsByUserID(ctx, db.ListExtractionsByUserIDParams{
		UserID:  userID,
		Column2: int64(limit),
		Column3: int64(offset),
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []Extraction{}, nil
	}

	gearRows, err := r.q.ListExtractionGearByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	gearByExtraction := make(map[string][]GearInfo)
	for _, g := range gearRows {
		gearByExtraction[g.ExtractionID] = append(gearByExtraction[g.ExtractionID], GearInfo{
			ID:     g.ID,
			TypeID: g.TypeID,
			Name:   g.Name,
		})
	}

	items := make([]Extraction, len(rows))
	for i, row := range rows {
		gear := gearByExtraction[row.ID]
		if gear == nil {
			gear = []GearInfo{}
		}
		items[i] = listRowToExtraction(row, gear)
	}
	return items, nil
}

// GetExtraction implements ExtractionRepo.
func (r *pgExtractionRepo) GetExtraction(ctx context.Context, id, userID string) (*Extraction, error) {
	row, err := r.q.GetExtractionByID(ctx, db.GetExtractionByIDParams{ID: id, UserID: userID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	gearRows, err := r.q.GetExtractionGear(ctx, id)
	if err != nil {
		return nil, err
	}
	gear := make([]GearInfo, len(gearRows))
	for i, g := range gearRows {
		gear[i] = GearInfo{ID: g.ID, TypeID: g.TypeID, Name: g.Name}
	}

	e := getRowToExtraction(row, gear)
	return &e, nil
}

// CreateExtraction implements ExtractionRepo.
func (r *pgExtractionRepo) CreateExtraction(ctx context.Context, userID string, p ExtractionParams) (*Extraction, error) {
	beanRow, err := r.q.GetBeanByID(ctx, db.GetBeanByIDParams{ID: p.BeanID, UserID: userID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidBean
		}
		return nil, err
	}

	gear, err := r.resolveGear(ctx, userID, p.GearIDs)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	qtx := r.q.WithTx(tx)

	row, err := qtx.CreateExtraction(ctx, db.CreateExtractionParams{
		UserID:      userID,
		BeanID:      p.BeanID,
		DoseIn:      p.DoseIn,
		YieldOut:    p.YieldOut,
		Time:        p.Time,
		TargetTime:  p.TargetTime,
		GrindSize:   p.GrindSize,
		PreInfusion: p.PreInfusion,
		TastingNote: toNullStr(p.TastingNote),
	})
	if err != nil {
		return nil, err
	}

	for _, gid := range p.GearIDs {
		if err := qtx.InsertExtractionGear(ctx, db.InsertExtractionGearParams{
			ExtractionID: row.ID,
			GearID:       gid,
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Extraction{
		ID:     row.ID,
		UserID: row.UserID,
		BeanID: row.BeanID,
		Bean: BeanInfo{
			ID:      beanRow.ID,
			Name:    beanRow.Name,
			Roaster: nullStr(beanRow.Roaster),
			Roast:   nullStr(beanRow.RoastLevel),
		},
		DoseIn:      row.DoseIn,
		YieldOut:    row.YieldOut,
		Time:        row.Time,
		TargetTime:  row.TargetTime,
		GrindSize:   row.GrindSize,
		PreInfusion: row.PreInfusion,
		TastingNote: nullStr(row.TastingNote),
		Gear:        gear,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

// UpdateExtraction implements ExtractionRepo.
func (r *pgExtractionRepo) UpdateExtraction(ctx context.Context, id, userID string, p ExtractionParams) (*Extraction, error) {
	beanRow, err := r.q.GetBeanByID(ctx, db.GetBeanByIDParams{ID: p.BeanID, UserID: userID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidBean
		}
		return nil, err
	}

	gear, err := r.resolveGear(ctx, userID, p.GearIDs)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	qtx := r.q.WithTx(tx)

	row, err := qtx.UpdateExtraction(ctx, db.UpdateExtractionParams{
		ID:          id,
		UserID:      userID,
		BeanID:      p.BeanID,
		DoseIn:      p.DoseIn,
		YieldOut:    p.YieldOut,
		Time:        p.Time,
		TargetTime:  p.TargetTime,
		GrindSize:   p.GrindSize,
		PreInfusion: p.PreInfusion,
		TastingNote: toNullStr(p.TastingNote),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := qtx.DeleteExtractionGearByExtractionID(ctx, id); err != nil {
		return nil, err
	}

	for _, gid := range p.GearIDs {
		if err := qtx.InsertExtractionGear(ctx, db.InsertExtractionGearParams{
			ExtractionID: id,
			GearID:       gid,
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Extraction{
		ID:     row.ID,
		UserID: row.UserID,
		BeanID: row.BeanID,
		Bean: BeanInfo{
			ID:      beanRow.ID,
			Name:    beanRow.Name,
			Roaster: nullStr(beanRow.Roaster),
			Roast:   nullStr(beanRow.RoastLevel),
		},
		DoseIn:      row.DoseIn,
		YieldOut:    row.YieldOut,
		Time:        row.Time,
		TargetTime:  row.TargetTime,
		GrindSize:   row.GrindSize,
		PreInfusion: row.PreInfusion,
		TastingNote: nullStr(row.TastingNote),
		Gear:        gear,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

// DeleteExtraction implements ExtractionRepo.
func (r *pgExtractionRepo) DeleteExtraction(ctx context.Context, id, userID string) error {
	n, err := r.q.DeleteExtractionByID(ctx, db.DeleteExtractionByIDParams{ID: id, UserID: userID})
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// resolveGear validates that all gearIDs belong to userID and returns ordered GearInfo items.
func (r *pgExtractionRepo) resolveGear(ctx context.Context, userID string, gearIDs []string) ([]GearInfo, error) {
	if len(gearIDs) == 0 {
		return []GearInfo{}, nil
	}
	allGear, err := r.q.ListGearByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	gearMap := make(map[string]db.Gear, len(allGear))
	for _, g := range allGear {
		gearMap[g.ID] = g
	}
	result := make([]GearInfo, len(gearIDs))
	for i, id := range gearIDs {
		g, ok := gearMap[id]
		if !ok {
			return nil, ErrInvalidGear
		}
		result[i] = GearInfo{ID: g.ID, TypeID: g.TypeID, Name: g.Name}
	}
	return result, nil
}

func getRowToExtraction(row db.GetExtractionByIDRow, gear []GearInfo) Extraction {
	return Extraction{
		ID:     row.ID,
		UserID: row.UserID,
		BeanID: row.BeanID,
		Bean: BeanInfo{
			ID:      row.BeanID,
			Name:    row.BeanName,
			Roaster: nullStr(row.BeanRoaster),
			Roast:   nullStr(row.BeanRoast),
		},
		DoseIn:      row.DoseIn,
		YieldOut:    row.YieldOut,
		Time:        row.Time,
		TargetTime:  row.TargetTime,
		GrindSize:   row.GrindSize,
		PreInfusion: row.PreInfusion,
		TastingNote: nullStr(row.TastingNote),
		Gear:        gear,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func listRowToExtraction(row db.ListExtractionsByUserIDRow, gear []GearInfo) Extraction {
	return Extraction{
		ID:     row.ID,
		UserID: row.UserID,
		BeanID: row.BeanID,
		Bean: BeanInfo{
			ID:      row.BeanID,
			Name:    row.BeanName,
			Roaster: nullStr(row.BeanRoaster),
			Roast:   nullStr(row.BeanRoast),
		},
		DoseIn:      row.DoseIn,
		YieldOut:    row.YieldOut,
		Time:        row.Time,
		TargetTime:  row.TargetTime,
		GrindSize:   row.GrindSize,
		PreInfusion: row.PreInfusion,
		TastingNote: nullStr(row.TastingNote),
		Gear:        gear,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func nullStr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	v := s.String
	return &v
}

func toNullStr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
