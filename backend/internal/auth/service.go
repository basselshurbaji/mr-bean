package auth

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/basselshurbaji/mr_bean/backend/internal/mailer"
	"github.com/basselshurbaji/mr_bean/backend/internal/user"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, rawRefreshToken string) (accessToken, refreshToken string, err error)
	Register(ctx context.Context, firstName, lastName, email, password string) (accessToken, refreshToken string, err error)
}

type authService struct {
	users  user.UserRepo
	tokens TokenService
	mailer mailer.Mailer
}

func NewAuthService(users user.UserRepo, tokens TokenService, mailer mailer.Mailer) AuthService {
	return &authService{users: users, tokens: tokens, mailer: mailer}
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

func (s *authService) Register(ctx context.Context, firstName, lastName, email, password string) (string, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	u, err := s.users.Create(ctx, firstName, lastName, email, string(hash))
	if err != nil {
		return "", "", errors.New("email already registered")
	}

	accessToken, err := s.tokens.GenerateAccessToken(u.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokens.GenerateRefreshToken(u.ID)
	if err != nil {
		return "", "", err
	}

	go func() {
		if err := s.mailer.Send(context.Background(), mailer.Email{
			To:       email,
			Subject:  "Welcome to Mr. Bean",
			Template: "welcome.html",
			Data:     mailer.WelcomeData{FirstName: firstName},
		}); err != nil {
			log.Printf("send welcome email to %s: %v", email, err)
		}
	}()

	return accessToken, refreshToken, nil
}
