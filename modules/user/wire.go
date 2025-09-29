//go:build wireinject
// +build wireinject

package user

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/winartodev/apollo-be/modules/user/delivery/http"
)

func InitializeUserAPI(
	db *sql.DB,
) (*http.UserHandler, error) {
	wire.Build(moduleSet)
	return &http.UserHandler{}, nil
}
