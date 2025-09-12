//go:build wireinject
// +build wireinject

package country

import (
	"github.com/google/wire"
	"github.com/winartodev/apollo-be/modules/country/delivery/http"
)

func InitializeCountryAPI() (*http.CountryHandler, error) {
	wire.Build(moduleSet)
	return &http.CountryHandler{}, nil
}
