package usecase

import (
	"context"
	"github.com/winartodev/apollo-be/internal/domain/entities"
	domainError "github.com/winartodev/apollo-be/internal/domain/error"

	"github.com/winartodev/apollo-be/modules/user/domain/service"
	"github.com/winartodev/apollo-be/modules/user/usecase/dto"
)

type UserUseCase interface {
	GetCurrentUser(ctx context.Context) (res *dto.UserDto, err error)
	CheckUserIfExists(ctx context.Context, data *entities.SharedUser) (res bool, err error)
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

func (uc *userUseCase) CheckUserIfExists(ctx context.Context, data *entities.SharedUser) (res bool, err error) {
	var exists bool

	exists, err = uc.userService.IsEmailExists(ctx, data.Email)
	if err != nil {
		return false, err
	}
	if exists {
		return true, domainError.ErrEmailAlreadyExists
	}

	exists, err = uc.userService.IsUsernameExists(ctx, data.Username)
	if err != nil {
		return false, err
	}
	if exists {
		return true, domainError.ErrUsernameAlreadyExists
	}

	return false, nil
}
