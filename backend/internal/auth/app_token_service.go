package auth

import (
	"context"
	"errors"
)

type AppTokenService interface {
	Create(ctx context.Context, userID, appName string) (token string, record *AppToken, err error)
	Revoke(ctx context.Context, tokenID, userID string) error
	Validate(ctx context.Context, rawToken string) (userID string, err error)
}

type appTokenService struct {
	repo   AppTokenRepo
	tokens TokenService
}

func NewAppTokenService(repo AppTokenRepo, tokens TokenService) AppTokenService {
	return &appTokenService{repo: repo, tokens: tokens}
}

// Create implements AppTokenService.
func (s *appTokenService) Create(ctx context.Context, userID, appName string) (string, *AppToken, error) {
	record, err := s.repo.Create(ctx, userID, appName)
	if err != nil {
		return "", nil, err
	}
	token, err := s.tokens.GenerateAppToken(userID, record.ID)
	if err != nil {
		return "", nil, err
	}
	return token, record, nil
}

// Revoke implements AppTokenService.
func (s *appTokenService) Revoke(ctx context.Context, tokenID, userID string) error {
	return s.repo.Revoke(ctx, tokenID, userID)
}

// Validate implements AppTokenService.
func (s *appTokenService) Validate(ctx context.Context, rawToken string) (string, error) {
	claims, err := s.tokens.ValidateAppToken(rawToken)
	if err != nil {
		return "", err
	}
	record, err := s.repo.GetByID(ctx, claims.ID)
	if err != nil {
		return "", errors.New("token not found")
	}
	if record.Revoked {
		return "", errors.New("token revoked")
	}
	return claims.UserID, nil
}