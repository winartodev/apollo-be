package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/winartodev/apollo-be/infrastructure/http/response"
	"github.com/winartodev/apollo-be/infrastructure/middleware"
	"github.com/winartodev/apollo-be/modules/user/usecase"
)

type UserHandler struct {
	middleware  *middleware.Middleware
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase, middleware *middleware.Middleware) *UserHandler {
	return &UserHandler{
		middleware:  middleware,
		userUseCase: userUseCase,
	}
}

func (uh *UserHandler) GetUserInfo(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := uh.userUseCase.GetCurrentUser(ctx)
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, http.StatusOK, "OK", res.ToResponse(), nil)
}

func (uh *UserHandler) RegisterRoutes(api *echo.Group) error {

	user := api.Group("/users")
	user.GET("/me", uh.GetUserInfo, uh.middleware.HandleWithAuth())

	return nil
}
