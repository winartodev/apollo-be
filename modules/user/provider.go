package user

import (
	"github.com/google/wire"
	"github.com/winartodev/apollo-be/infrastructure/provider"
	"github.com/winartodev/apollo-be/modules/user/delivery/http"
	userService "github.com/winartodev/apollo-be/modules/user/domain/service"
	userRepo "github.com/winartodev/apollo-be/modules/user/repository"
	userUseCase "github.com/winartodev/apollo-be/modules/user/usecase"
)

var repositorySet = wire.NewSet(
	// Repository implementations
	userRepo.NewUserRepository,
)

var serviceSet = wire.NewSet(
	// Domain services
	userService.NewUserService,
)

var useCaseSet = wire.NewSet(
	// Use cases
	userUseCase.NewUserUseCase,
)

var handlerSet = wire.NewSet(
	// HTTP Handlers
	http.NewUserHandler,
)

var moduleSet = wire.NewSet(
	provider.InfraProviderSet,
	provider.MiddlewareProviderSet,
	repositorySet,
	serviceSet,
	useCaseSet,
	handlerSet,
)
