package country

import (
	"github.com/google/wire"
	"github.com/winartodev/apollo-be/modules/country/delivery/http"
)

var handlerSet = wire.NewSet(
	http.NewCountryHandler,
)

var moduleSet = wire.NewSet(
	handlerSet,
)
