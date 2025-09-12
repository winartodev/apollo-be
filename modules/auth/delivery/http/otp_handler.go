package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/winartodev/apollo-be/infrastructure/http/response"
	"github.com/winartodev/apollo-be/infrastructure/middleware"
	"github.com/winartodev/apollo-be/modules/auth/delivery/http/dto"
	"github.com/winartodev/apollo-be/modules/auth/usecase"
)

type OtpHandler struct {
	middleware *middleware.Middleware
	otpUseCase usecase.OtpUseCase
}

func NewOtpHandler(otpUseCase usecase.OtpUseCase, middleware *middleware.Middleware) *OtpHandler {
	return &OtpHandler{
		middleware: middleware,
		otpUseCase: otpUseCase,
	}
}

// ResendOtp handles OTP resend requests
//
//	@Summary		Resend OTP
//	@Description	Resend one-time password to the user
//	@Tags			OTP
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.OtpResendRequest							true	"OTP resend request data"
//	@Success		200		{object}	response.Response{data=dto.OtpRefreshResponse}	"OTP resent successfully"
//	@Failure		400		{object}	response.ErrorResponse							"Invalid request payload"
//	@Failure		401		{object}	response.ErrorResponse							"Unauthorized"
//	@Failure		422		{object}	response.ErrorResponse							"Validation error"
//	@Failure		429		{object}	response.ErrorResponse							"Too many requests - rate limited"
//	@Failure		500		{object}	response.ErrorResponse							"Internal server error"
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
	res, err := oh.otpUseCase.ResendOTP(ctx)
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	resp := dto.OtpRefreshResponse{
		RetryAttemptsLeft: res.RetryAttemptsLeft,
		RetryAfterIn:      res.RetryAfterIn,
		ExpiresIn:         res.ExpiresIn,
		IsValid:           res.IsValid,
	}

	return response.SuccessResponse(c, http.StatusOK, "ok", resp, nil)
}

// ValidateOtp handles OTP validation requests
//
//	@Summary		Validate OTP
//	@Description	Validate one-time password provided by the user
//	@Tags			OTP
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
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

	ctx := c.Request().Context()
	res, err := oh.otpUseCase.ValidateOTP(ctx, fmt.Sprintf("%d", req.OTPNumber))
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	resp := dto.OtpValidationResponse{
		IsValid:         res.IsValid,
		RedirectionLink: "/homePage",
	}

	return response.SuccessResponse(c, http.StatusOK, "ok", resp, nil)
}

func (oh *OtpHandler) RegisterRoutes(api *echo.Group) error {

	otp := api.Group("/otp")
	otp.POST("/resend", oh.ResendOtp, oh.middleware.HandleWithAuth())
	otp.POST("/validate", oh.ValidateOtp, oh.middleware.HandleWithAuth())
	return nil
}
