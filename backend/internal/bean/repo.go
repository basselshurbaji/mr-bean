package bean

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "github.com/basselshurbaji/mr_bean/backend/db/sqlc"
)

var ErrNotFound = errors.New("not found")

type Bean struct {
	ID           string
	UserID       string
	Name         string
	Roaster      *string
	Origin       *string
	Process      *string
	RoastLevel   *string
	TastingNotes *string
	Notes        *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type BeanParams struct {
	Name         string
	Roaster      *string
	Origin       *string
	Process      *string
	RoastLevel   *string
	TastingNotes *string
	Notes        *string
}

type BeanRepo interface {
	ListBeans(ctx context.Context, userID string) ([]Bean, error)
	CreateBean(ctx context.Context, userID string, p BeanParams) (*Bean, error)
	UpdateBean(ctx context.Context, id, userID string, p BeanParams) (*Bean, error)
	DeleteBean(ctx context.Context, id, userID string) error
}

type pgBeanRepo struct {
	q *db.Queries
}

func NewPgBeanRepo(d *sql.DB) BeanRepo {
	return &pgBeanRepo{q: db.New(d)}
}

// ListBeans implements BeanRepo.
func (r *pgBeanRepo) ListBeans(ctx context.Context, userID string) ([]Bean, error) {
	rows, err := r.q.ListBeansByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	items := make([]Bean, len(rows))
	for i, row := range rows {
		items[i] = rowToBean(row)
	}
	return items, nil
}

// CreateBean implements BeanRepo.
func (r *pgBeanRepo) CreateBean(ctx context.Context, userID string, p BeanParams) (*Bean, error) {
	row, err := r.q.CreateBean(ctx, db.CreateBeanParams{
		UserID:       userID,
		Name:         p.Name,
		Roaster:      toNullStr(p.Roaster),
		Origin:       toNullStr(p.Origin),
		Process:      toNullStr(p.Process),
		RoastLevel:   toNullStr(p.RoastLevel),
		TastingNotes: toNullStr(p.TastingNotes),
		Notes:        toNullStr(p.Notes),
	})
	if err != nil {
		return nil, err
	}
	b := rowToBean(row)
	return &b, nil
}

// UpdateBean implements BeanRepo.
func (r *pgBeanRepo) UpdateBean(ctx context.Context, id, userID string, p BeanParams) (*Bean, error) {
	row, err := r.q.UpdateBean(ctx, db.UpdateBeanParams{
		ID:           id,
		UserID:       userID,
		Name:         p.Name,
		Roaster:      toNullStr(p.Roaster),
		Origin:       toNullStr(p.Origin),
		Process:      toNullStr(p.Process),
		RoastLevel:   toNullStr(p.RoastLevel),
		TastingNotes: toNullStr(p.TastingNotes),
		Notes:        toNullStr(p.Notes),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	b := rowToBean(row)
	return &b, nil
}

// DeleteBean implements BeanRepo.
func (r *pgBeanRepo) DeleteBean(ctx context.Context, id, userID string) error {
	n, err := r.q.DeleteBeanByID(ctx, db.DeleteBeanByIDParams{ID: id, UserID: userID})
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func rowToBean(row db.Bean) Bean {
	return Bean{
		ID:           row.ID,
		UserID:       row.UserID,
		Name:         row.Name,
		Roaster:      nullStr(row.Roaster),
		Origin:       nullStr(row.Origin),
		Process:      nullStr(row.Process),
		RoastLevel:   nullStr(row.RoastLevel),
		TastingNotes: nullStr(row.TastingNotes),
		Notes:        nullStr(row.Notes),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
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
