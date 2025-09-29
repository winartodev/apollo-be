package service

import (
	"context"

	"github.com/winartodev/apollo-be/internal/domain"
	"github.com/winartodev/apollo-be/internal/domain/entities"
	domainError "github.com/winartodev/apollo-be/internal/domain/error"
	"github.com/winartodev/apollo-be/modules/auth/domain/repository"
)

type AuthService interface {
	CreateNewUser(ctx context.Context, data entities.SharedUser) (res *entities.SharedUser, err error)
	VerifyUser(ctx context.Context, username string, password string) (res *entities.SharedUser, err error)
	UpdateRefreshToken(ctx context.Context, id int64, token *string) (err error)
}

type authService struct {
	passwordService domain.PasswordService
	authRepo        repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository, passwordService domain.PasswordService) (AuthService, error) {
	return &authService{
		passwordService: passwordService,
		authRepo:        authRepo,
	}, nil
}

func (as *authService) CreateNewUser(ctx context.Context, data entities.SharedUser) (res *entities.SharedUser, err error) {
	encryptedPassword, err := as.passwordService.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	data.Password = encryptedPassword

	id, err := as.authRepo.RegisterNewUserDB(ctx, data)
	if err != nil {
		return nil, err
	}

	data.ID = *id

	return &data, nil
}

func (as *authService) VerifyUser(ctx context.Context, username string, password string) (res *entities.SharedUser, err error) {
	user, err := as.authRepo.GetUserDataDB(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domainError.ErrUserNotFound
	}

	if !as.passwordService.ComparePassword(password, user.Password) {
		return nil, domainError.ErrInvalidUsernameOrPassword
	}

	return user, nil
}

func (as *authService) UpdateRefreshToken(ctx context.Context, id int64, token *string) (err error) {
	return as.authRepo.UpdateRefreshTokenDB(ctx, id, token)
}
