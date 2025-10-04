package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/winartodev/apollo-be/config"
	domainError "github.com/winartodev/apollo-be/internal/domain/error"
	"github.com/winartodev/apollo-be/modules/auth/domain/entities"
	"github.com/winartodev/apollo-be/modules/auth/domain/repository"
	"math/big"
	"time"
)

type OtpService interface {
	GetOTP(ctx context.Context, username string) (service *entities.OTPService, err error)
	ValidateOTP(ctx context.Context, username string, otp *string) (valid bool, err error)
}

type otpService struct {
	otpRepo repository.OtpRepository
	otp     *config.Otp
}

func NewOtpService(otpRepo repository.OtpRepository, otp *config.Otp) (OtpService, error) {
	return &otpService{
		otpRepo: otpRepo,
		otp:     otp,
	}, nil
}

func (os *otpService) GetOTP(ctx context.Context, username string) (service *entities.OTPService, err error) {
	otp, err := os.generateOTP(6)
	if err != nil {
		return nil, err
	}

	currentAttempt, err := os.otpRepo.GetOtpAttemptRedis(ctx, username)
	if err != nil {
		return nil, err
	}

	if currentAttempt != nil && *currentAttempt >= os.otp.MaxAttempt {
		return nil, domainError.ErrOtpTooManyRequest
	}

	exp := time.Duration(os.otp.Expiration) * time.Second
	err = os.otpRepo.SetOtpRedis(ctx, username, entities.OTP{
		Number: *otp,
	}, exp)
	if err != nil {
		return nil, err
	}

	incr, err := os.otpRepo.IncrOtpAttemptRedis(ctx, username)
	if err != nil {
		return nil, err
	}

	return &entities.OTPService{
		Expiration:        os.otp.Expiration,
		MaxAttempt:        os.otp.MaxAttempt,
		RetryInterval:     os.otp.RetryInterval,
		RetryAttemptsLeft: os.otp.MaxAttempt - *incr,
		OTPNumber:         *otp,
	}, nil
}

func (os *otpService) ValidateOTP(ctx context.Context, username string, otp *string) (valid bool, err error) {
	otpData, err := os.otpRepo.GetOtpRedis(ctx, username)
	if err != nil {
		return false, err
	}

	if otpData == nil {
		return false, domainError.ErrInvalidOTPNumber
	}

	if otpData.Number != *otp {
		return false, domainError.ErrInvalidOTPNumber
	}

	return true, err
}

func (os *otpService) generateOTP(length int) (res *string, err error) {
	if length <= 0 {
		return nil, fmt.Errorf("length must be positive")
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number: %v", err)
		}
		result[i] = byte(num.Int64() + '0')
	}

	resultStr := string(result)

	return &resultStr, nil
}
