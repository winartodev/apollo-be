package middleware

import (
	"context"
	"net/http"
	"strings"

	domainError "github.com/winartodev/apollo-be/internal/domain/error"

	"github.com/labstack/echo/v4"
	customContext "github.com/winartodev/apollo-be/infrastructure/context"
	"github.com/winartodev/apollo-be/infrastructure/http/response"
	"github.com/winartodev/apollo-be/internal/domain"
)

type Middleware struct {
	jwt domain.TokenService
}

func NewMiddleware(jwt domain.TokenService) *Middleware {
	return &Middleware{jwt: jwt}
}

func (m *Middleware) HandleWithAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("Authorization")

			claims, err := m.verifyToken(c.Request().Context(), authorization, true)
			if err != nil {
				return response.FailedResponse(c, http.StatusUnauthorized, err)
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, customContext.UserIdKey, claims.UserID)
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set(string(customContext.UserIdKey), claims.UserID)

			return next(c)
		}
	}
}

func (m *Middleware) HandleRefreshToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("Authorization")

			claims, err := m.verifyToken(c.Request().Context(), authorization, false)
			if err != nil {
				return response.FailedResponse(c, http.StatusUnauthorized, err)
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, customContext.UserIdKey, claims.UserID)
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set(string(customContext.UserIdKey), claims.UserID)

			return next(c)
		}
	}
}

func (m *Middleware) HandleWithAPIKey() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("X-API-Key")

			if m.verifyAPIKey(apiKey) {
				return response.FailedResponse(c, http.StatusUnauthorized, domainError.ErrInvalidAuthorizationHeader)
			}

			return next(c)
		}
	}
}

func (m *Middleware) verifyToken(ctx context.Context, authorization string, isAccessToken bool) (claims *domain.TokenClaims, err error) {
	if authorization == "" {
		return nil, domainError.ErrAuthorizationHeaderEmpty
	}

	parts := strings.Split(authorization, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, domainError.ErrInvalidAuthorizationHeader
	}

	token := parts[1]
	if token == "" {
		return nil, domainError.ErrEmptyToken
	}

	if isAccessToken {
		claims, err = m.jwt.ValidateAccessToken(token)
	} else {
		claims, err = m.jwt.ValidateRefreshToken(token)
	}

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (m *Middleware) verifyAPIKey(apiKey string) (isValid bool) {
	return "hello" == apiKey
}

func GetAppPlatform() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			platform := c.Request().Header.Get("X-APP-PLATFORM")
			if platform != "" {
				c.Set(string(customContext.AppPlatformKey), platform)
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, customContext.AppPlatformKey, platform)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
