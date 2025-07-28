//go:build wireinject
// +build wireinject

package user

import (
	"github.com/google/wire"
	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/modules/user/delivery/http"
)

func InitializeUserAPI(
	database *helper.DatabaseUtil,
	redis *helper.RedisUtil,
	jwt *helper.JWT,
) (*http.UserHandler, error) {
	wire.Build(moduleSet)
	return &http.UserHandler{}, nil
}
