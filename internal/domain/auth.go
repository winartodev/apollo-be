package domain

import (
	"context"
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	UserID    int64
	Username  string
	Email     string
	IssueAt   time.Time
	ExpiresAt time.Time
}

type TokenService interface {
	GenerateTokenPair(ctx context.Context, user *SharedUser) (*TokenPair, error)
	ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error)
	ValidateRefreshToken(ctx context.Context, token string) (*TokenClaims, error)
	InvalidateToken(ctx context.Context, token string) error
}

type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) bool
}
