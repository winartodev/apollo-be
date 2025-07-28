//go:build wireinject
// +build wireinject

package auth

import (
	"github.com/google/wire"
	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/modules/auth/delivery/http"
)

func InitializeAuthAPI(
	database *helper.DatabaseUtil,
	redis *helper.RedisUtil,
	jwt *helper.JWT,
	smtp *helper.SmtpConfig,
	otp *helper.OtpConfig,
) (*http.AuthHandler, error) {
	wire.Build(moduleSet)
	return &http.AuthHandler{}, nil
}

func InitializeOtpAPI(
	database *helper.DatabaseUtil,
	redis *helper.RedisUtil,
	jwt *helper.JWT,
	smtp *helper.SmtpConfig,
	otp *helper.OtpConfig,
) (*http.OtpHandler, error) {
	wire.Build(moduleSet)
	return &http.OtpHandler{}, nil
}
