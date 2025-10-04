package domain

import (
	"time"

	domainEntity "github.com/winartodev/apollo-be/internal/domain/entities"
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
	GenerateTokenPair(user *domainEntity.SharedUser) (*TokenPair, error)
	ValidateAccessToken(token string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*TokenClaims, error)
	InvalidateToken(token string) error
}

type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) bool
}
