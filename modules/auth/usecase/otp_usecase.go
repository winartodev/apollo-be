package usecase

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"

	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/modules/auth/domain/service"
	"github.com/winartodev/apollo-be/modules/auth/usecase/dto"
	userUseCase "github.com/winartodev/apollo-be/modules/user/usecase"
)

type OtpUseCase interface {
	ResendOTP(ctx context.Context) (res *dto.OtpDto, err error)
	ValidateOTP(ctx context.Context, code string) (res *dto.OtpDto, err error)
}

type otpUseCase struct {
	jwt         *helper.JWT
	smtp        *helper.SmtpConfig
	otp         *helper.OtpConfig
	userUseCase userUseCase.UserUseCase
	otpService  service.OtpService
}

func NewOtpUseCase(otpService service.OtpService, userUseCase userUseCase.UserUseCase, jwt *helper.JWT, smtp *helper.SmtpConfig, otp *helper.OtpConfig) OtpUseCase {
	return &otpUseCase{
		jwt:         jwt,
		smtp:        smtp,
		otp:         otp,
		otpService:  otpService,
		userUseCase: userUseCase,
	}
}

func (ou *otpUseCase) ResendOTP(ctx context.Context) (res *dto.OtpDto, err error) {
	user, err := ou.userUseCase.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	otp, retryLeft, err := ou.otpService.GetOTP(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	err = ou.sendOTPEmail(user.Email, *otp)
	if err != nil {
		return nil, err
	}

	retryAttemptsLeft := ou.otp.MaxAttempt - *retryLeft

	return &dto.OtpDto{
		ExpiresIn:         ou.otp.Expiration,
		RetryAfterIn:      ou.otp.Expiration,
		RetryAttemptsLeft: retryAttemptsLeft,
		IsValid:           false,
	}, nil
}

func (ou *otpUseCase) ValidateOTP(ctx context.Context, code string) (res *dto.OtpDto, err error) {
	user, err := ou.userUseCase.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	otpIsValid, err := ou.otpService.ValidateOTP(ctx, user.Email, &code)
	if err != nil {
		return nil, err
	}

	return &dto.OtpDto{
		IsValid: otpIsValid,
	}, nil
}

func (ou *otpUseCase) sendOTPEmail(email string, code string) (err error) {
	_, filename, _, _ := runtime.Caller(0)
	templatePath := filepath.Join(filepath.Dir(filename), "templates", "otp_email_template.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	data := make(map[string]interface{})
	data["otp"] = code
	data["exp"] = 3

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to render email template: %v", err)
	}

	err = ou.smtp.SendSmtpHTML(email, "Your Verification Code", body.String())
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %v", err)
	}

	return nil
}
