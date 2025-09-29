package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/winartodev/apollo-be/infrastructure/http/response"
)

func HealthCheck(c echo.Context) error {
	return response.SuccessResponse(c, http.StatusOK, "OK", nil, nil)
}
