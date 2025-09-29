package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/winartodev/apollo-be/helper"
	infraContext "github.com/winartodev/apollo-be/infrastructure/context"
	"github.com/winartodev/apollo-be/infrastructure/http/response"
	"github.com/winartodev/apollo-be/infrastructure/middleware"
	domainError "github.com/winartodev/apollo-be/internal/domain/error"
	"github.com/winartodev/apollo-be/modules/auth/delivery/enums"
	"github.com/winartodev/apollo-be/modules/auth/delivery/http/dto"
	"github.com/winartodev/apollo-be/modules/auth/usecase"
)

type AuthHandler struct {
	middleware  *middleware.Middleware
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, middleware *middleware.Middleware) *AuthHandler {
	return &AuthHandler{
		middleware:  middleware,
		authUseCase: authUseCase,
	}
}

// SignUp handles user registration
//
//	@Summary		Register a new user
//	@Description	Create a new user account with provided credentials
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SignUpRequest							true	"User registration data"
//	@Success		201		{object}	response.Response{data=dto.AuthResponse}	"User registered successfully"
//	@Failure		400		{object}	response.ErrorResponse						"Invalid request payload"
//	@Failure		422		{object}	response.ErrorResponse						"Validation error"
//	@Failure		500		{object}	response.ErrorResponse						"Internal server error"
//	@Router			/auth/sign-up [post]
func (ah *AuthHandler) SignUp(c echo.Context) error {
	var req dto.SignUpRequest

	if err := c.Bind(&req); err != nil {
		return response.FailedResponse(c, http.StatusBadRequest, fmt.Errorf(response.ErrInvalidRequestPayload, err))
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrResponse(c, err)
	}

	ctx := c.Request().Context()
	res, err := ah.authUseCase.SignUp(ctx, req.ToUseCaseData())
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	resp := dto.AuthResponse{
		AccessToken:     res.AccessToken,
		RefreshToken:    res.RefreshToken,
		RedirectionLink: ah.buildRedirectionLink(ctx, enums.SignUp),
		Otp: &dto.OtpRefreshResponse{
			RetryAttemptsLeft: res.Otp.RetryAttemptsLeft,
			ExpiresIn:         res.Otp.ExpiresIn,
			RetryAfterIn:      res.Otp.RetryAfterIn,
		},
	}

	return response.SuccessResponse(c, http.StatusCreated, "User registered successfully", resp, nil)
}

// SignIn handles user authentication
//
//	@Summary		Authenticate user
//	@Description	Sign in user with username/email and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SignInRequest							true	"User credentials"
//	@Success		200		{object}	response.Response{data=dto.AuthResponse}	"User authenticated successfully"
//	@Failure		400		{object}	response.ErrorResponse						"Invalid request payload"
//	@Failure		401		{object}	response.ErrorResponse						"Invalid username or password"
//	@Failure		422		{object}	response.ErrorResponse						"Validation error"
//	@Failure		500		{object}	response.ErrorResponse						"Internal server error"
//	@Router			/auth/sign-in [post]
func (ah *AuthHandler) SignIn(c echo.Context) error {
	var req dto.SignInRequest

	if err := c.Bind(&req); err != nil {
		return response.FailedResponse(c, http.StatusBadRequest, fmt.Errorf(response.ErrInvalidRequestPayload, err))
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrResponse(c, err)
	}

	ctx := c.Request().Context()
	res, err := ah.authUseCase.SignIn(ctx, req.ToUseCaseData())
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	resp := dto.AuthResponse{
		AccessToken:     res.AccessToken,
		RefreshToken:    res.RefreshToken,
		RedirectionLink: ah.buildRedirectionLink(ctx, enums.SignIn),
	}

	return response.SuccessResponse(c, http.StatusOK, "OK", resp, nil)
}

// SignOut handles user logout
//
//	@Summary		Logout user
//	@Description	Sign out user and invalidate authentication tokens
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=dto.AuthResponse}	"User signed out successfully"
//	@Failure		401	{object}	response.ErrorResponse						"Unauthorized - Invalid or missing token"
//	@Failure		500	{object}	response.ErrorResponse						"Internal server error"
//	@Router			/auth/sign-out [post]
func (ah *AuthHandler) SignOut(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := ah.authUseCase.SignOut(ctx)
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	resp := dto.AuthResponse{
		AccessToken:     res.AccessToken,
		RefreshToken:    res.RefreshToken,
		RedirectionLink: ah.buildRedirectionLink(ctx, enums.SignOut),
	}

	return response.SuccessResponse(c, http.StatusOK, "OK", resp, nil)
}

// RefreshToken handles token refresh
//
//	@Summary		Refresh authentication tokens
//	@Description	Refresh access token using refresh token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=dto.AuthResponse}	"Tokens refreshed successfully"
//	@Failure		401	{object}	response.ErrorResponse						"Unauthorized - Invalid or expired refresh token"
//	@Failure		500	{object}	response.ErrorResponse						"Internal server error"
//	@Router			/auth/refresh [post]
func (ah *AuthHandler) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := ah.authUseCase.RefreshToken(ctx)
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	resp := dto.AuthResponse{
		AccessToken:     res.AccessToken,
		RefreshToken:    res.RefreshToken,
		RedirectionLink: "/",
	}

	return response.SuccessResponse(c, http.StatusOK, "OK", resp, nil)
}

func (ah *AuthHandler) VerifyUser(c echo.Context) error {
	var req dto.VerifyUserRequest

	if err := c.Bind(&req); err != nil {
		return response.FailedResponse(c, http.StatusBadRequest, fmt.Errorf(response.ErrInvalidRequestPayload, err))
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrResponse(c, err)
	}

	ctx := c.Request().Context()
	err := ah.authUseCase.VerifyUser(ctx, req.Username)
	if err != nil {
		return response.FailedResponse(c, http.StatusConflict, domainError.ErrInvalidUsernameOrPassword)
	}

	return response.SuccessResponse(c, http.StatusOK, "ok", nil, nil)
}

func (ah *AuthHandler) RegisterRoutes(api *echo.Group) error {
	auth := api.Group("/auth")
	auth.POST("/sign-up", ah.SignUp)
	auth.POST("/sign-in", ah.SignIn)
	auth.GET("/verify-user", ah.VerifyUser)
	auth.POST("/sign-out", ah.SignOut, ah.middleware.HandleWithAuth())
	auth.POST("/refresh", ah.RefreshToken, ah.middleware.HandleRefreshToken())

	return nil
}

func (ah *AuthHandler) buildRedirectionLink(ctx context.Context, action enums.AuthOperation) string {
	platform, err := infraContext.GetAppPlatformFromContext(ctx)
	if err != nil {
		return ah.getWebRedirection(action)
	}

	return helper.BuildRedirectionLink[enums.AuthOperation](
		platform,
		action,
		ah.getMobileRedirection,
		ah.getWebRedirection,
	)
}

func (ah *AuthHandler) getMobileRedirection(operation enums.AuthOperation) (res string) {
	switch operation {
	case enums.SignUp:
		return "/otpVerificationPage"
	case enums.SignIn:
		return "/homePage"
	case enums.SignOut:
		return "/SignInPage"
	default:
		return ""
	}
}

func (ah *AuthHandler) getWebRedirection(operation enums.AuthOperation) (res string) {
	switch operation {
	case enums.SignUp:
		return "/verification"
	case enums.SignIn:
		return "/home"
	case enums.SignOut:
		return "/sign-in"
	default:
		return ""
	}
}
