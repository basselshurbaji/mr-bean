package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// UserStore is the minimal interface auth needs for credential validation.
type UserStore interface {
	GetByEmail(ctx context.Context, email string) (*StoredUser, error)
}

// StoredUser carries only what auth needs — ID, hashed password, active flag.
type StoredUser struct {
	ID           string
	PasswordHash string
	IsActive     bool
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, rawRefreshToken string) (accessToken, refreshToken string, err error)
}

type authService struct {
	users  UserStore
	tokens TokenService
}

func NewAuthService(users UserStore, tokens TokenService) AuthService {
	return &authService{users: users, tokens: tokens}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !u.IsActive {
		return "", "", errors.New("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := s.tokens.GenerateAccessToken(u.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokens.GenerateRefreshToken(u.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Refresh(ctx context.Context, rawRefreshToken string) (string, string, error) {
	claims, err := s.tokens.ValidateRefreshToken(rawRefreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	accessToken, err := s.tokens.GenerateAccessToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokens.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
