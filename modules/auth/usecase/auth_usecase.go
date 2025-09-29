package usecase

import (
	"context"

	infraContext "github.com/winartodev/apollo-be/infrastructure/context"
	"github.com/winartodev/apollo-be/internal/domain"
	"github.com/winartodev/apollo-be/internal/domain/entities"
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
	VerifyUser(ctx context.Context, username string) (err error)
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
	var sharedUser = &entities.SharedUser{
		Username:    data.Username,
		Password:    data.Password,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
	}

	userExists, err := uc.userUseCase.CheckUserIfExists(ctx, sharedUser)
	if err != nil {
		return nil, err
	}

	if userExists {
		return nil, domainError.ErrUserAlreadyExists
	}

	newUser, err := uc.authService.CreateNewUser(ctx, *sharedUser)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, infraContext.UserIdKey, newUser.ID)
	otp, err := uc.otpUseCase.ResendOTP(ctx)
	if err != nil {
		return nil, err
	}

	domainSharedUser := &domain.SharedUser{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
	}

	jwt, err := uc.jwt.GenerateTokenPair(ctx, domainSharedUser)
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

	sharedUser := &domain.SharedUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	jwt, err := uc.jwt.GenerateTokenPair(ctx, sharedUser)
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

	sharedUser := &domain.SharedUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	jwt, err := uc.jwt.GenerateTokenPair(ctx, sharedUser)
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
