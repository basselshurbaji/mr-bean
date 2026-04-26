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
	Create(ctx context.Context, firstName, lastName, email, passwordHash string) (*User, error)
	UpdateProfile(ctx context.Context, id, firstName, lastName string) (*User, error)
	UpdatePassword(ctx context.Context, id, passwordHash string) error
}

type pgUserRepo struct {
	q *db.Queries
}

func NewPgUserRepo(d *sql.DB) UserRepo {
	return &pgUserRepo{q: db.New(d)}
}

// GetByEmail implements UserRepo.
func (r *pgUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return rowToUser(row), nil
}

// GetByID implements UserRepo.
func (r *pgUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return rowToUser(row), nil
}

// Create implements UserRepo.
func (r *pgUserRepo) Create(ctx context.Context, firstName, lastName, email, passwordHash string) (*User, error) {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, err
	}
	return rowToUser(row), nil
}

// UpdateProfile implements UserRepo.
func (r *pgUserRepo) UpdateProfile(ctx context.Context, id, firstName, lastName string) (*User, error) {
	row, err := r.q.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		return nil, err
	}
	return rowToUser(row), nil
}

// UpdatePassword implements UserRepo.
func (r *pgUserRepo) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	return r.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
	})
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
