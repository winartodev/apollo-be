package auth

import (
	"fmt"
	"time"

	domainEntity "github.com/winartodev/apollo-be/internal/domain/entities"

	"github.com/winartodev/apollo-be/internal/domain"
)

type JwtTokenService struct {
	jwt *JWT
}

func NewJwtTokenService(jwt *JWT) domain.TokenService {
	return &JwtTokenService{
		jwt: jwt,
	}
}

// GenerateTokenPair implements domain.TokenService.
func (jts *JwtTokenService) GenerateTokenPair(user *domainEntity.SharedUser) (*domain.TokenPair, error) {
	userJWT := &UserJWT{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	tokenPair, err := jts.jwt.GenerateToken(userJWT)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %v", err)
	}

	return &domain.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

// InvalidateToken implements domain.TokenService.
func (jts *JwtTokenService) InvalidateToken(token string) error {
	// TODO: Implement invalidate token
	return nil
}

// ValidateAccessToken implements domain.TokenService.
func (jts *JwtTokenService) ValidateAccessToken(token string) (*domain.TokenClaims, error) {
	claims, isValid, err := jts.jwt.VerifyToken(jts.jwt.AccessToken.SecretKey, token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify access token: %v", err)
	}

	if !isValid {
		return nil, fmt.Errorf("invalid access token")
	}

	return jts.claimsToTokenClaims(claims)
}

// ValidateRefreshToken implements domain.TokenService.
func (jts *JwtTokenService) ValidateRefreshToken(token string) (*domain.TokenClaims, error) {
	claims, isValid, err := jts.jwt.VerifyToken(jts.jwt.RefreshToken.SecretKey, token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify refresh token: %v", err)
	}

	if !isValid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return jts.claimsToTokenClaims(claims)
}

func (jts *JwtTokenService) claimsToTokenClaims(claims map[string]interface{}) (*domain.TokenClaims, error) {
	tokenClaims := &domain.TokenClaims{}
	if id, ok := claims["id"].(float64); ok {
		tokenClaims.UserID = int64(id)
	}

	if username, ok := claims["username"].(string); ok {
		tokenClaims.Username = username
	}

	if email, ok := claims["email"].(string); ok {
		tokenClaims.Email = email
	}

	if issueAt, ok := claims["issueAt"].(float64); ok {
		tokenClaims.IssueAt = time.Unix(int64(issueAt), 0)
	}

	if expireAt, ok := claims["expireAt"].(float64); ok {
		tokenClaims.ExpiresAt = time.Unix(int64(expireAt), 0)
	}

	return tokenClaims, nil
}
