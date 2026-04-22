package user

import (
	"context"
	"database/sql"
	"time"

	db "github.com/basselshurbaji/mr_bean/backend/db/sqlc"
)

type User struct {
	ID           string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type pgUserRepo struct {
	q *db.Queries
}

func NewPgUserRepo(d *sql.DB) UserRepo {
	return &pgUserRepo{q: db.New(d)}
}

func (r *pgUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return rowToUser(row), nil
}

func (r *pgUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return rowToUser(row), nil
}

func rowToUser(row db.User) *User {
	return &User{
		ID:           row.ID,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
