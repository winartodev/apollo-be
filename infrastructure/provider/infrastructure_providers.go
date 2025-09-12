package provider

import (
	"github.com/google/wire"
	"github.com/winartodev/apollo-be/infrastructure/auth"
	"github.com/winartodev/apollo-be/infrastructure/database"
	"github.com/winartodev/apollo-be/infrastructure/middleware"
	"github.com/winartodev/apollo-be/infrastructure/redis"
	"github.com/winartodev/apollo-be/infrastructure/smtp"
	"github.com/winartodev/apollo-be/internal/application/service"
)

// InfraProviderSet contains infrastructure implementations
var InfraProviderSet = wire.NewSet(
	// Infrastructure services
	auth.NewJWT,
	auth.NewJwtTokenService,
	auth.NewBcryptPasswordService,
	database.NewDatabase,
	redis.NewRedis,
	smtp.NewSMTPService,
)

// MiddlewareProviderSet contains middleware components
var MiddlewareProviderSet = wire.NewSet(
	middleware.NewMiddleware,
)

// ApplicationServiceProviderSet contains application services
var ApplicationServiceProviderSet = wire.NewSet(
	service.NewUserApplicationService,
)
