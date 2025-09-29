package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	domainError "github.com/winartodev/apollo-be/internal/domain/error"
	"github.com/winartodev/apollo-be/modules/auth/domain/entities"
	"github.com/winartodev/apollo-be/modules/auth/domain/repository"
)

const (
	otpExp         = 3 * time.Minute
	otpMaxAttempts = 3
)

type OtpService interface {
	GetOTP(ctx context.Context, username string) (otp *string, retryLeft *int64, err error)
	ValidateOTP(ctx context.Context, username string, otp *string) (valid bool, err error)
}

type otpService struct {
	otpRepo repository.OtpRepository
}

func NewOtpService(otpRepo repository.OtpRepository) (OtpService, error) {
	return &otpService{
		otpRepo: otpRepo,
	}, nil
}

func (os *otpService) GetOTP(ctx context.Context, username string) (otp *string, retryLeft *int64, err error) {
	otp, err = os.generateOTP(6)
	if err != nil {
		return nil, nil, err
	}

	currentAttempt, err := os.otpRepo.GetOtpAttemptRedis(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	if currentAttempt != nil && *currentAttempt >= otpMaxAttempts {

		return nil, nil, domainError.ErrOtpTooManyRequest
	}

	err = os.otpRepo.SetOtpRedis(ctx, username, entities.OTP{
		Number: *otp,
	}, otpExp)
	if err != nil {
		return nil, nil, err
	}

	incr, err := os.otpRepo.IncrOtpAttemptRedis(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	return otp, incr, nil
}

func (os *otpService) ValidateOTP(ctx context.Context, username string, otp *string) (valid bool, err error) {
	otpData, err := os.otpRepo.GetOtpRedis(ctx, username)
	if err != nil {
		return false, err
	}

	if otpData == nil {
		return false, nil
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
