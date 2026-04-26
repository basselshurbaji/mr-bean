package user

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*User, error)
	UpdateProfile(ctx context.Context, id, firstName, lastName string) (*User, error)
	ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error
}

type userService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) UserService {
	return &userService{repo: repo}
}

// GetByID implements UserService.
func (s *userService) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateProfile implements UserService.
func (s *userService) UpdateProfile(ctx context.Context, id, firstName, lastName string) (*User, error) {
	if firstName == "" || lastName == "" {
		u, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if firstName == "" {
			firstName = u.FirstName
		}
		if lastName == "" {
			lastName = u.LastName
		}
	}
	return s.repo.UpdateProfile(ctx, id, firstName, lastName)
}

// ChangePassword implements UserService.
func (s *userService) ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("incorrect password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, id, string(hash))
}
