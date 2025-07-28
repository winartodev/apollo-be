package usecase

import (
	"context"

	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/modules/user/domain/service"
	"github.com/winartodev/apollo-be/modules/user/usecase/dto"
	userUseCaseError "github.com/winartodev/apollo-be/modules/user/usecase/error"
)

type UserUseCase interface {
	GetCurrentUser(ctx context.Context) (res *dto.UserDto, err error)
	CheckUserIfExists(ctx context.Context, username string) (res bool, err error)
}

type userUseCase struct {
	userService service.UserService
}

func NewUserUseCase(userService service.UserService) (UserUseCase, error) {
	return &userUseCase{
		userService: userService,
	}, nil
}

func (uc *userUseCase) GetCurrentUser(ctx context.Context) (res *dto.UserDto, err error) {
	user, err := uc.userService.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	userDto := user.ToUseCaseData()

	return &userDto, nil
}

func (uc *userUseCase) CheckUserIfExists(ctx context.Context, username string) (res bool, err error) {
	var exists bool

	if helper.IsEmailValid(username) {
		exists, err = uc.userService.IsEmailExists(ctx, username)
		if err != nil {
			return false, err
		}
		if exists {
			return true, userUseCaseError.EmailAlreadyExists
		}
	} else {
		exists, err = uc.userService.IsUsernameExists(ctx, username)
		if err != nil {
			return false, err
		}
		if exists {
			return true, userUseCaseError.UsernameAlreadyExists
		}
	}

	return false, nil
}
