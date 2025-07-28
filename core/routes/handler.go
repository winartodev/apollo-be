package routes

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	errorRegisterHandler = "failed to register handler: %v"
)

type APIRouteItf interface {
	RegisterRoutes(api *echo.Group) error
}

func RegisterHandler(e *echo.Echo, handlers ...APIRouteItf) error {
	api := e.Group("/api")
	api.GET("/health-check", HealthCheck)

	for _, apiRoute := range handlers {
		if err := apiRoute.RegisterRoutes(api); err != nil {
			return errors.New(fmt.Sprintf(errorRegisterHandler, err))
		}
	}

	return nil
}
