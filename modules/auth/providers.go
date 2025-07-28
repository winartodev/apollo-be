package auth

import (
	"github.com/google/wire"
	authHttp "github.com/winartodev/apollo-be/modules/auth/delivery/http"
	authService "github.com/winartodev/apollo-be/modules/auth/domain/service"
	authRepo "github.com/winartodev/apollo-be/modules/auth/repository"
	authUsecase "github.com/winartodev/apollo-be/modules/auth/usecase"
	userService "github.com/winartodev/apollo-be/modules/user/domain/service"
	userRepo "github.com/winartodev/apollo-be/modules/user/repository"
	userUseCase "github.com/winartodev/apollo-be/modules/user/usecase"
)

var repositorySet = wire.NewSet(
	// Repository implementations
	authRepo.NewAuthRepository,
	authRepo.NewOtpRepository,
	userRepo.NewUserRepository,
)

var serviceSet = wire.NewSet(
	// Domain services
	authService.NewAuthService,
	authService.NewOtpService,
	userService.NewUserService,
)

var useCaseSet = wire.NewSet(
	// Use cases
	authUsecase.NewAuthUseCase,
	authUsecase.NewOtpUseCase,
	userUseCase.NewUserUseCase,
)

var handlerSet = wire.NewSet(
	// HTTP Handlers
	authHttp.NewAuthHandler,
	authHttp.NewOtpHandler,
)

var moduleSet = wire.NewSet(
	repositorySet,
	serviceSet,
	useCaseSet,
	handlerSet,
)
