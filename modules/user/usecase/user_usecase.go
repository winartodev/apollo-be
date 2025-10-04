package usecase

import (
	"context"
	"github.com/winartodev/apollo-be/helper"

	domainEntity "github.com/winartodev/apollo-be/internal/domain/entities"
	domainError "github.com/winartodev/apollo-be/internal/domain/error"

	"github.com/winartodev/apollo-be/modules/user/domain/service"
	"github.com/winartodev/apollo-be/modules/user/usecase/dto"
)

type UserUseCase interface {
	GetCurrentUser(ctx context.Context) (res *dto.UserDto, err error)
	GetUserByEmail(ctx context.Context, email string) (res *domainEntity.SharedUser, err error)
	CheckUserIfExists(ctx context.Context, data domainEntity.SharedUser) (res *domainEntity.SharedUser, err error)
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

	if (user == nil) || (user.ID == 0) {
		return nil, domainError.ErrUserNotFound
	}

	userDto := user.ToUseCaseData()

	return &userDto, nil
}

func (uc *userUseCase) CheckUserIfExists(ctx context.Context, data domainEntity.SharedUser) (res *domainEntity.SharedUser, err error) {
	var sharedUser *domainEntity.SharedUser

	if data.Email != "" {
		sharedUser, err = uc.userService.IsEmailExists(ctx, data.Email)
		if err != nil {
			return nil, err
		}
		if sharedUser != nil {
			return sharedUser, domainError.ErrEmailAlreadyExists
		}
	}

	if data.Username != "" {
		sharedUser, err = uc.userService.IsUsernameExists(ctx, data.Username)
		if err != nil {
			return nil, err
		}
		if sharedUser != nil {
			return sharedUser, domainError.ErrUsernameAlreadyExists
		}
	}

	return nil, nil
}

func (uc *userUseCase) GetUserByEmail(ctx context.Context, email string) (res *domainEntity.SharedUser, err error) {
	if !helper.IsEmailValid(email) {
		return nil, domainError.ErrInvalidEmail
	}

	user, err := uc.userService.IsEmailExists(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domainError.ErrUserNotFound
	}

	return user, err
}
