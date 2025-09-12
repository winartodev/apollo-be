//go:build wireinject
// +build wireinject

package auth

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	config2 "github.com/winartodev/apollo-be/config"
	"github.com/winartodev/apollo-be/modules/auth/delivery/http"
)

func InitializeAuthAPI(
	db *sql.DB,
	redis *redis.Client,
	smtpConfig *config2.SMTPConfig,
	otp *config2.Otp,
) (*http.AuthHandler, error) {
	wire.Build(moduleSet)
	return &http.AuthHandler{}, nil
}

func InitializeOtpAPI(
	db *sql.DB,
	redis *redis.Client,
	smtpConfig *config2.SMTPConfig,
	otp *config2.Otp,
) (*http.OtpHandler, error) {
	wire.Build(moduleSet)
	return &http.OtpHandler{}, nil
}
