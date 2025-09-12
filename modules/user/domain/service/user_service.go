package service

import (
	"context"

	"github.com/winartodev/apollo-be/helper"
	infraContext "github.com/winartodev/apollo-be/infrastructure/context"
	domainError "github.com/winartodev/apollo-be/internal/domain/error"
	"github.com/winartodev/apollo-be/modules/user/domain/entities"
	"github.com/winartodev/apollo-be/modules/user/domain/repository"
)

type UserService interface {
	GetCurrentUser(ctx context.Context) (user *entities.User, err error)
	IsEmailExists(ctx context.Context, email string) (res bool, err error)
	IsUsernameExists(ctx context.Context, username string) (res bool, err error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) (UserService, error) {
	return &userService{userRepo: userRepo}, nil
}

func (us *userService) GetCurrentUser(ctx context.Context) (user *entities.User, err error) {
	id, err := infraContext.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err = us.userRepo.GetUserByIDDB(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) IsEmailExists(ctx context.Context, email string) (res bool, err error) {
	if !helper.IsEmailValid(email) {
		return false, domainError.ErrInvalidEmail
	}

	user, err := us.userRepo.GetUserByEmailDB(ctx, email)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func (us *userService) IsUsernameExists(ctx context.Context, username string) (res bool, err error) {
	user, err := us.userRepo.GetUserByUsernameDB(ctx, username)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}
