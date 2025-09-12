package service

import (
	"context"

	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/modules/user/domain/entities"
	"github.com/winartodev/apollo-be/modules/user/domain/errors"
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
	id, err := helper.GetUserIDFromContext(ctx)
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
		return false, errors.InvalidEmail
	}

	user, err := us.userRepo.GetUserByEmailDB(ctx, email)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func (us *userService) IsUsernameExists(ctx context.Context, username string) (res bool, err error) {
	user, err := us.userRepo.GetUserByEmailDB(ctx, username)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}
