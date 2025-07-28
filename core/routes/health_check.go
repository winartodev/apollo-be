package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/winartodev/apollo-be/core/helper"
)

func HealthCheck(c echo.Context) error {
	return helper.SuccessResponse(c, http.StatusOK, "OK", nil, nil)
}
