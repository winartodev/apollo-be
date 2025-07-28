package service

import (
	"context"

	"github.com/winartodev/apollo-be/modules/auth/domain/entities"
	"github.com/winartodev/apollo-be/modules/auth/domain/repository"
	"golang.org/x/crypto/bcrypt"

	domainError "github.com/winartodev/apollo-be/modules/auth/domain/error"
)

type AuthService interface {
	CreateNewUser(ctx context.Context, data entities.User) (res *entities.User, err error)
	VerifyUser(ctx context.Context, username string, password string) (res *entities.User, err error)
	UpdateRefreshToken(ctx context.Context, id int64, token *string) (err error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) (AuthService, error) {
	return &authService{
		authRepo: authRepo,
	}, nil
}

func (as *authService) CreateNewUser(ctx context.Context, data entities.User) (res *entities.User, err error) {
	encryptedPassword, err := as.hashPassword(data.Password)
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

func (as *authService) VerifyUser(ctx context.Context, username string, password string) (res *entities.User, err error) {
	user, err := as.authRepo.GetUserDataDB(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domainError.ErrUserNotFound
	}

	if !as.verifyPassword(password, user.Password) {
		return nil, domainError.ErrInvalidUsernameOrPassword
	}

	return user, nil
}

func (as *authService) UpdateRefreshToken(ctx context.Context, id int64, token *string) (err error) {
	return as.authRepo.UpdateRefreshTokenDB(ctx, id, token)
}

func (as *authService) hashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (as *authService) verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
