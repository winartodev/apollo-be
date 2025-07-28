package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/core/middleware"
	"github.com/winartodev/apollo-be/modules/user/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (uh *UserHandler) GetUserInfo(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := uh.userUseCase.GetCurrentUser(ctx)
	if err != nil {
		return helper.FailedResponse(c, http.StatusInternalServerError, err)
	}

	return helper.SuccessResponse(c, http.StatusOK, "OK", res.ToResponse(), nil)
}

func (uh *UserHandler) RegisterRoutes(api *echo.Group) error {

	user := api.Group("/users")
	user.GET("/me", uh.GetUserInfo, middleware.HandleWithAuth())

	return nil
}
