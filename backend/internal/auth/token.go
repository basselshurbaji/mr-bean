package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID    string    `json:"uid"`
	TokenType TokenType `json:"typ"`
}

type TokenService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateAccessToken(raw string) (*Claims, error)
	ValidateRefreshToken(raw string) (*Claims, error)
}

type jwtService struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewTokenService(secret string, accessExpiry, refreshExpiry time.Duration) TokenService {
	return &jwtService{
		secret:        []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateAccessToken implements TokenService.
func (s *jwtService) GenerateAccessToken(userID string) (string, error) {
	return s.generate(userID, TokenTypeAccess, s.accessExpiry)
}

// GenerateRefreshToken implements TokenService.
func (s *jwtService) GenerateRefreshToken(userID string) (string, error) {
	return s.generate(userID, TokenTypeRefresh, s.refreshExpiry)
}

// ValidateAccessToken implements TokenService.
func (s *jwtService) ValidateAccessToken(raw string) (*Claims, error) {
	return s.validate(raw, TokenTypeAccess)
}

// ValidateRefreshToken implements TokenService.
func (s *jwtService) ValidateRefreshToken(raw string) (*Claims, error) {
	return s.validate(raw, TokenTypeRefresh)
}

func (s *jwtService) generate(userID string, tokenType TokenType, expiry time.Duration) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:    userID,
		TokenType: tokenType,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *jwtService) validate(raw string, expectedType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenType != expectedType {
		return nil, errors.New("wrong token type")
	}
	return claims, nil
}
