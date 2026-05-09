package auth

import (
	"context"
	"database/sql"
	"time"

	db "github.com/basselshurbaji/mr_bean/backend/db/sqlc"
)

type AppToken struct {
	ID        string
	UserID    string
	AppName   string
	Revoked   bool
	CreatedAt time.Time
}

type AppTokenRepo interface {
	Create(ctx context.Context, userID, appName string) (*AppToken, error)
	GetByID(ctx context.Context, id string) (*AppToken, error)
	Revoke(ctx context.Context, id, userID string) error
}

type pgAppTokenRepo struct {
	q *db.Queries
}

func NewPgAppTokenRepo(d *sql.DB) AppTokenRepo {
	return &pgAppTokenRepo{q: db.New(d)}
}

// Create implements AppTokenRepo.
func (r *pgAppTokenRepo) Create(ctx context.Context, userID, appName string) (*AppToken, error) {
	row, err := r.q.CreateAppToken(ctx, db.CreateAppTokenParams{
		UserID:  userID,
		AppName: appName,
	})
	if err != nil {
		return nil, err
	}
	return rowToAppToken(row), nil
}

// GetByID implements AppTokenRepo.
func (r *pgAppTokenRepo) GetByID(ctx context.Context, id string) (*AppToken, error) {
	row, err := r.q.GetAppTokenByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return rowToAppToken(row), nil
}

// Revoke implements AppTokenRepo.
func (r *pgAppTokenRepo) Revoke(ctx context.Context, id, userID string) error {
	return r.q.RevokeAppToken(ctx, db.RevokeAppTokenParams{
		ID:     id,
		UserID: userID,
	})
}

func rowToAppToken(row db.AppToken) *AppToken {
	return &AppToken{
		ID:        row.ID,
		UserID:    row.UserID,
		AppName:   row.AppName,
		Revoked:   row.Revoked,
		CreatedAt: row.CreatedAt,
	}
}