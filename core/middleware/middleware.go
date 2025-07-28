package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/winartodev/apollo-be/core/helper"
)

var (
	globalJWT *helper.JWT
)

func NewMiddleware(jwt *helper.JWT) {
	globalJWT = jwt
}

func GetAppPlatform() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			platform := c.Request().Header.Get("X-APP-PLATFORM")
			if platform != "" {
				c.Set(helper.AppPlatformKey, platform)
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, helper.AppPlatformKey, platform)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func HandleWithAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("Authorization")

			claims, err := verifyToken(authorization, globalJWT.AccessToken.SecretKey)
			if err != nil {
				return helper.FailedResponse(c, http.StatusUnauthorized, err)
			}

			userId := int64(claims["id"].(float64))
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, helper.UserIdKey, userId)
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set(helper.UserIdKey, userId)

			return next(c)
		}
	}
}

func HandleRefreshToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("Authorization")

			claims, err := verifyToken(authorization, globalJWT.RefreshToken.SecretKey)
			if err != nil {
				return helper.FailedResponse(c, http.StatusUnauthorized, err)
			}

			userId := int64(claims["id"].(float64))
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, helper.UserIdKey, userId)
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set(helper.UserIdKey, userId)

			return next(c)
		}
	}
}

func verifyToken(authorization string, secretKey []byte) (claims map[string]interface{}, err error) {
	fmt.Println(string(secretKey))
	if authorization == "" {
		return nil, errorAuthorizationHeaderEmpty
	}

	parts := strings.Split(authorization, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errorInvalidAuthorizationHeader
	}

	token := parts[1]
	if token == "" {
		return nil, errorEmptyToken
	}

	claims, isValid, err := globalJWT.VerifyToken(secretKey, token)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, errInvalidToken
	}

	return claims, nil
}
