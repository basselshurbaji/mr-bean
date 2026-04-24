package gear

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "github.com/basselshurbaji/mr_bean/backend/db/sqlc"
)

var ErrNotFound = errors.New("not found")
var ErrUnownedGear = errors.New("one or more gear IDs are unknown or do not belong to you")

type GearItem struct {
	ID        string
	UserID    string
	TypeID    string
	Name      string
	Brand     *string
	Model     *string
	Year      *string
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateGearParams struct {
	TypeID string
	Name   string
	Brand  *string
	Model  *string
	Year   *string
	Notes  *string
}

type UpdateGearParams = CreateGearParams

type Station struct {
	ID        string
	UserID    string
	Name      string
	Gear      []GearItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GearRepo interface {
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

type pgGearRepo struct {
	db *sql.DB
	q  *db.Queries
}

func NewPgGearRepo(d *sql.DB) GearRepo {
	return &pgGearRepo{db: d, q: db.New(d)}
}

func (r *pgGearRepo) ListGear(ctx context.Context, userID string) ([]GearItem, error) {
	rows, err := r.q.ListGearByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	items := make([]GearItem, len(rows))
	for i, row := range rows {
		items[i] = rowToGearItem(row)
	}
	return items, nil
}

func (r *pgGearRepo) GetGear(ctx context.Context, id, userID string) (*GearItem, error) {
	row, err := r.q.GetGearByID(ctx, db.GetGearByIDParams{ID: id, UserID: userID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	g := rowToGearItem(row)
	return &g, nil
}

func (r *pgGearRepo) CreateGear(ctx context.Context, userID string, p CreateGearParams) (*GearItem, error) {
	row, err := r.q.CreateGear(ctx, db.CreateGearParams{
		UserID: userID,
		TypeID: p.TypeID,
		Name:   p.Name,
		Brand:  toNullStr(p.Brand),
		Model:  toNullStr(p.Model),
		Year:   toNullStr(p.Year),
		Notes:  toNullStr(p.Notes),
	})
	if err != nil {
		return nil, err
	}
	g := rowToGearItem(row)
	return &g, nil
}

func (r *pgGearRepo) UpdateGear(ctx context.Context, id, userID string, p UpdateGearParams) (*GearItem, error) {
	row, err := r.q.UpdateGear(ctx, db.UpdateGearParams{
		ID:     id,
		UserID: userID,
		TypeID: p.TypeID,
		Name:   p.Name,
		Brand:  toNullStr(p.Brand),
		Model:  toNullStr(p.Model),
		Year:   toNullStr(p.Year),
		Notes:  toNullStr(p.Notes),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	g := rowToGearItem(row)
	return &g, nil
}

// DeleteGear removes a gear item and re-sequences station positions in a transaction.
func (r *pgGearRepo) DeleteGear(ctx context.Context, id, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	qtx := r.q.WithTx(tx)

	stationIDs, err := qtx.GetStationIDsByGearID(ctx, id)
	if err != nil {
		return err
	}

	for _, sid := range stationIDs {
		remaining, err := qtx.ListGearIDsInStationExcluding(ctx, db.ListGearIDsInStationExcludingParams{
			StationID: sid,
			GearID:    id,
		})
		if err != nil {
			return err
		}
		if err := qtx.DeleteStationGearByStationID(ctx, sid); err != nil {
			return err
		}
		for pos, gid := range remaining {
			if err := qtx.InsertStationGear(ctx, db.InsertStationGearParams{
				StationID: sid,
				GearID:    gid,
				Position:  int32(pos),
			}); err != nil {
				return err
			}
		}
	}

	n, err := qtx.DeleteGearByID(ctx, db.DeleteGearByIDParams{ID: id, UserID: userID})
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}

func (r *pgGearRepo) ListStations(ctx context.Context, userID string) ([]Station, error) {
	stationRows, err := r.q.ListStationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(stationRows) == 0 {
		return []Station{}, nil
	}

	gearRows, err := r.q.ListStationGearByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	gearByStation := make(map[string][]GearItem)
	for _, row := range gearRows {
		gearByStation[row.StationID] = append(gearByStation[row.StationID], joinRowToGearItem(row))
	}

	stations := make([]Station, len(stationRows))
	for i, s := range stationRows {
		gear := gearByStation[s.ID]
		if gear == nil {
			gear = []GearItem{}
		}
		stations[i] = Station{
			ID:        s.ID,
			UserID:    s.UserID,
			Name:      s.Name,
			Gear:      gear,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		}
	}
	return stations, nil
}

// CreateStation creates the station record and inserts gear associations in a transaction.
func (r *pgGearRepo) CreateStation(ctx context.Context, userID, name string, gearIDs []string) (*Station, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	qtx := r.q.WithTx(tx)

	s, err := qtx.CreateStation(ctx, db.CreateStationParams{UserID: userID, Name: name})
	if err != nil {
		return nil, err
	}

	for pos, gid := range gearIDs {
		if err := qtx.InsertStationGear(ctx, db.InsertStationGearParams{
			StationID: s.ID,
			GearID:    gid,
			Position:  int32(pos),
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Station{
		ID:        s.ID,
		UserID:    s.UserID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}, nil
}

// UpdateStation replaces the station name and gear list atomically.
func (r *pgGearRepo) UpdateStation(ctx context.Context, id, userID, name string, gearIDs []string) (*Station, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	qtx := r.q.WithTx(tx)

	s, err := qtx.UpdateStation(ctx, db.UpdateStationParams{ID: id, UserID: userID, Name: name})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := qtx.DeleteStationGearByStationID(ctx, id); err != nil {
		return nil, err
	}

	for pos, gid := range gearIDs {
		if err := qtx.InsertStationGear(ctx, db.InsertStationGearParams{
			StationID: id,
			GearID:    gid,
			Position:  int32(pos),
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Station{
		ID:        s.ID,
		UserID:    s.UserID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}, nil
}

func (r *pgGearRepo) DeleteStation(ctx context.Context, id, userID string) error {
	n, err := r.q.DeleteStationByID(ctx, db.DeleteStationByIDParams{ID: id, UserID: userID})
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func rowToGearItem(row db.Gear) GearItem {
	return GearItem{
		ID:        row.ID,
		UserID:    row.UserID,
		TypeID:    row.TypeID,
		Name:      row.Name,
		Brand:     nullStr(row.Brand),
		Model:     nullStr(row.Model),
		Year:      nullStr(row.Year),
		Notes:     nullStr(row.Notes),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func joinRowToGearItem(row db.ListStationGearByUserIDRow) GearItem {
	return GearItem{
		ID:        row.ID,
		UserID:    row.UserID,
		TypeID:    row.TypeID,
		Name:      row.Name,
		Brand:     nullStr(row.Brand),
		Model:     nullStr(row.Model),
		Year:      nullStr(row.Year),
		Notes:     nullStr(row.Notes),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
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
