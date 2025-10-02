package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/winartodev/apollo-be/helper"
	"github.com/winartodev/apollo-be/modules/auth/delivery/enums"

	infraContext "github.com/winartodev/apollo-be/infrastructure/context"
	userUseCase "github.com/winartodev/apollo-be/modules/user/usecase"

	"github.com/labstack/echo/v4"
	"github.com/winartodev/apollo-be/infrastructure/http/response"
	"github.com/winartodev/apollo-be/infrastructure/middleware"
	"github.com/winartodev/apollo-be/modules/auth/delivery/http/dto"
	"github.com/winartodev/apollo-be/modules/auth/usecase"
)

type OtpHandler struct {
	middleware  *middleware.Middleware
	otpUseCase  usecase.OtpUseCase
	userUseCase userUseCase.UserUseCase
}

func NewOtpHandler(otpUseCase usecase.OtpUseCase, userUseCase userUseCase.UserUseCase, middleware *middleware.Middleware) *OtpHandler {
	return &OtpHandler{
		middleware:  middleware,
		otpUseCase:  otpUseCase,
		userUseCase: userUseCase,
	}
}

// ResendOtp godoc
//
//	@Summary		Resend OTP
//	@Description	Resend one-time password to the user
//	@Tags			OTP
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.OtpResendRequest					true	"OTP resend request data"
//	@Success		200		{object}	response.Response{data=dto.OtpResponse}	"OTP resent successfully"
//	@Failure		400		{object}	response.ErrorResponse					"Invalid request payload"
//	@Failure		401		{object}	response.ErrorResponse					"Unauthorized"
//	@Failure		422		{object}	response.ErrorResponse					"Validation error"
//	@Failure		429		{object}	response.ErrorResponse					"Too many requests - rate limited"
//	@Failure		500		{object}	response.ErrorResponse					"Internal server error"
//	@Router			/otp/resend [post]
func (oh *OtpHandler) ResendOtp(c echo.Context) error {
	var req dto.OtpResendRequest

	if err := c.Bind(&req); err != nil {
		return response.FailedResponse(c, http.StatusBadRequest, fmt.Errorf(response.ErrInvalidRequestPayload, err))
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrResponse(c, err)
	}

	ctx := c.Request().Context()
	user, err := oh.userUseCase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	ctx = context.WithValue(ctx, infraContext.UserIdKey, user.ID)
	res, err := oh.otpUseCase.SendOTP(ctx)
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	resp := dto.OtpResponse{
		RetryAttemptsLeft: res.RetryAttemptsLeft,
		RetryAfterIn:      res.RetryAfterIn,
		ExpiresIn:         res.ExpiresIn,
		IsValid:           res.IsValid,
	}

	return response.SuccessResponse(c, http.StatusOK, "ok", resp, nil)
}

// ValidateOtp godoc
//
//	@Summary		Validate OTP
//	@Description	Validate one-time password provided by the user
//	@Tags			OTP
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.OtpRequest										true	"OTP validation data"
//	@Success		200		{object}	response.Response{data=dto.OtpValidationResponse}	"OTP validated successfully"
//	@Failure		400		{object}	response.ErrorResponse								"Invalid request payload"
//	@Failure		401		{object}	response.ErrorResponse								"Unauthorized"
//	@Failure		403		{object}	response.ErrorResponse								"Invalid OTP"
//	@Failure		422		{object}	response.ErrorResponse								"Validation error"
//	@Failure		429		{object}	response.ErrorResponse								"Too many attempts"
//	@Failure		500		{object}	response.ErrorResponse								"Internal server error"
//	@Router			/otp/validate [post]
func (oh *OtpHandler) ValidateOtp(c echo.Context) error {
	var req dto.OtpRequest

	if err := c.Bind(&req); err != nil {
		return response.FailedResponse(c, http.StatusBadRequest, fmt.Errorf(response.ErrInvalidRequestPayload, err))
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrResponse(c, err)
	}

	actionType, err := enums.ParseOtpOperationEnum(req.Type)
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	ctx := c.Request().Context()
	user, err := oh.userUseCase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	ctx = context.WithValue(ctx, infraContext.UserIdKey, user.ID)
	res, err := oh.otpUseCase.ValidateOTP(ctx, req.OTPNumber)
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	resp := dto.OtpValidationResponse{
		IsValid:         res.IsValid,
		RedirectionLink: oh.buildRedirectionLink(ctx, actionType),
	}

	return response.SuccessResponse(c, http.StatusOK, "ok", resp, nil)
}

func (oh *OtpHandler) RegisterRoutes(api *echo.Group) error {

	otp := api.Group("/otp")
	otp.POST("/resend", oh.ResendOtp)
	otp.POST("/validate", oh.ValidateOtp)

	return nil
}

func (oh *OtpHandler) buildRedirectionLink(ctx context.Context, action enums.OtpOperationEnum) string {
	platform, err := infraContext.GetAppPlatformFromContext(ctx)
	if err != nil {
		return oh.getWebRedirection(action)
	}

	return helper.BuildRedirectionLink[enums.OtpOperationEnum](
		platform,
		action,
		oh.getMobileRedirection,
		oh.getWebRedirection,
	)
}

func (oh *OtpHandler) getMobileRedirection(action enums.OtpOperationEnum) (res string) {
	switch action {
	case enums.OtpSignUp:
		return "/signInPage"
	case enums.OtpRequestReset:
		return "/resetPasswordPage"
	default:
		return ""
	}
}

func (oh *OtpHandler) getWebRedirection(action enums.OtpOperationEnum) (res string) {
	switch action {
	case enums.OtpSignUp:
		return "/sign-in"
	case enums.OtpRequestReset:
		return "/reset-password"
	default:
		return ""
	}
}
