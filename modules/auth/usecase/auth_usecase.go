package usecase

import (
	"context"
	"errors"
	"github.com/winartodev/apollo-be/helper"

	infraContext "github.com/winartodev/apollo-be/infrastructure/context"
	"github.com/winartodev/apollo-be/internal/domain"
	domainEntity "github.com/winartodev/apollo-be/internal/domain/entities"
	domainError "github.com/winartodev/apollo-be/internal/domain/error"
	authService "github.com/winartodev/apollo-be/modules/auth/domain/service"
	"github.com/winartodev/apollo-be/modules/auth/usecase/dto"
	userUseCase "github.com/winartodev/apollo-be/modules/user/usecase"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, data dto.SignUpDto) (res *dto.AuthDto, err error)
	SignIn(ctx context.Context, data dto.SignInDto) (res *dto.AuthDto, err error)
	SignOut(ctx context.Context) (res *dto.AuthDto, err error)
	RefreshToken(ctx context.Context) (res *dto.AuthDto, err error)
	VerifyUser(ctx context.Context, username string) (res *dto.VerifyUserDto, err error)
	RequestResetPassword(ctx context.Context, email string) (res *dto.AuthDto, err error)
	ResetPassword(ctx context.Context, data dto.ResetPasswordDto) (err error)
}

type authUseCase struct {
	jwt         domain.TokenService
	userUseCase userUseCase.UserUseCase
	authService authService.AuthService
	otpUseCase  OtpUseCase
}

func NewAuthUseCase(authService authService.AuthService, otpUseCase OtpUseCase, jwt domain.TokenService, userUseCase userUseCase.UserUseCase) (AuthUseCase, error) {
	return &authUseCase{
		jwt:         jwt,
		userUseCase: userUseCase,
		authService: authService,
		otpUseCase:  otpUseCase,
	}, nil
}

func (uc *authUseCase) SignUp(ctx context.Context, data dto.SignUpDto) (res *dto.AuthDto, err error) {
	var sharedUser = &domainEntity.SharedUser{
		Username:    data.Username,
		Password:    data.Password,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
	}

	user, err := uc.userUseCase.CheckUserIfExists(ctx, *sharedUser)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, domainError.ErrUserAlreadyExists
	}

	newUser, err := uc.authService.CreateNewUser(ctx, *sharedUser)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, infraContext.UserIdKey, newUser.ID)
	otp, err := uc.otpUseCase.SendOTP(ctx)
	if err != nil {
		return nil, err
	}

	domainSharedUser := &domainEntity.SharedUser{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
	}

	jwt, err := uc.jwt.GenerateTokenPair(domainSharedUser)
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
	user, err := uc.authService.VerifyUsernameAndPassword(ctx, data.Username, data.Password)
	if err != nil {
		return nil, err
	}

	sharedUser := &domainEntity.SharedUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	jwt, err := uc.jwt.GenerateTokenPair(sharedUser)
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
	id, err := infraContext.GetUserIDFromContext(ctx)
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

	sharedUser := &domainEntity.SharedUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	jwt, err := uc.jwt.GenerateTokenPair(sharedUser)
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

func (uc *authUseCase) VerifyUser(ctx context.Context, username string) (res *dto.VerifyUserDto, err error) {
	var sharedUser domainEntity.SharedUser

	if helper.IsEmailValid(username) {
		sharedUser.Email = username
	} else {
		sharedUser.Username = username
	}

	user, err := uc.userUseCase.CheckUserIfExists(ctx, sharedUser)
	if err != nil && (!errors.Is(err, domainError.ErrEmailAlreadyExists) && !errors.Is(err, domainError.ErrUsernameAlreadyExists)) {
		return nil, err
	}

	return &dto.VerifyUserDto{
		User: user,
	}, nil
}

func (uc *authUseCase) RequestResetPassword(ctx context.Context, email string) (res *dto.AuthDto, err error) {
	if !helper.IsEmailValid(email) {
		return nil, domainError.ErrInvalidEmail
	}

	user, err := uc.userUseCase.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, infraContext.UserIdKey, user.ID)
	otp, err := uc.otpUseCase.SendOTP(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.AuthDto{
		Otp: otp,
	}, nil
}

func (uc *authUseCase) ResetPassword(ctx context.Context, data dto.ResetPasswordDto) (err error) {
	if !uc.comparePassword(data.Password, data.PasswordConfirmation) {
		return domainError.ErrPasswordConfirmationMismatch
	}

	if !helper.IsEmailValid(data.Email) {
		return domainError.ErrInvalidEmail
	}

	user, err := uc.userUseCase.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return err
	}

	return uc.authService.UpdatePassword(ctx, user.ID, data.Password)
}

func (uc *authUseCase) comparePassword(password string, passwordConfirmation string) bool {
	return password == passwordConfirmation
}
