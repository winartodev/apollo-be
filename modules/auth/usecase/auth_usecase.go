package usecase

import (
	"context"

	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/modules/auth/domain/entities"
	authDomainError "github.com/winartodev/apollo-be/modules/auth/domain/error"
	authService "github.com/winartodev/apollo-be/modules/auth/domain/service"
	"github.com/winartodev/apollo-be/modules/auth/usecase/dto"
	userUseCase "github.com/winartodev/apollo-be/modules/user/usecase"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, data dto.SignUpDto) (res *dto.AuthDto, err error)
	SignIn(ctx context.Context, data dto.SignInDto) (res *dto.AuthDto, err error)
	SignOut(ctx context.Context) (res *dto.AuthDto, err error)
	RefreshToken(ctx context.Context) (res *dto.AuthDto, err error)
	VerifyUser(ctx context.Context, username string) (err error)
}

type authUseCase struct {
	jwt         *helper.JWT
	userUseCase userUseCase.UserUseCase
	authService authService.AuthService
	otpUseCase  OtpUseCase
}

func NewAuthUseCase(authService authService.AuthService, otpUseCase OtpUseCase, jwt *helper.JWT, userUseCase userUseCase.UserUseCase) (AuthUseCase, error) {
	return &authUseCase{
		jwt:         jwt,
		userUseCase: userUseCase,
		authService: authService,
		otpUseCase:  otpUseCase,
	}, nil
}

func (uc *authUseCase) SignUp(ctx context.Context, data dto.SignUpDto) (res *dto.AuthDto, err error) {
	userExists, err := uc.userUseCase.CheckUserIfExists(ctx, data.Username)
	if err != nil {
		return nil, err
	}

	if userExists {
		return nil, authDomainError.ErrUserAlreadyExists
	}

	newUser, err := uc.authService.CreateNewUser(ctx, entities.User{
		Username:    data.Username,
		Password:    data.Password,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
	})

	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, helper.UserIdKey, newUser.ID)
	otp, err := uc.otpUseCase.ResendOTP(ctx)
	if err != nil {
		return nil, err
	}

	jwt, err := uc.jwt.GenerateToken(&helper.UserJWT{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
	})
	if err != nil {
		return nil, err
	}

	err = uc.authService.UpdateRefreshToken(ctx, newUser.ID, &jwt.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.AuthDto{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		Otp:          otp,
	}, nil
}

func (uc *authUseCase) SignIn(ctx context.Context, data dto.SignInDto) (res *dto.AuthDto, err error) {
	user, err := uc.authService.VerifyUser(ctx, data.Username, data.Password)
	if err != nil {
		return nil, err
	}

	jwt, err := uc.jwt.GenerateToken(&helper.UserJWT{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return nil, err
	}

	err = uc.authService.UpdateRefreshToken(ctx, user.ID, &jwt.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.AuthDto{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
	}, nil
}

func (uc *authUseCase) SignOut(ctx context.Context) (res *dto.AuthDto, err error) {
	id, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = uc.authService.UpdateRefreshToken(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	return &dto.AuthDto{}, nil
}

func (uc *authUseCase) RefreshToken(ctx context.Context) (res *dto.AuthDto, err error) {
	user, err := uc.userUseCase.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	jwt, err := uc.jwt.GenerateToken(&helper.UserJWT{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return nil, err
	}

	err = uc.authService.UpdateRefreshToken(ctx, user.ID, &jwt.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.AuthDto{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
	}, nil
}

func (uc *authUseCase) VerifyUser(ctx context.Context, username string) (err error) {

	return nil
}
