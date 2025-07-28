package repository

import (
	"context"
	"time"

	"github.com/winartodev/apollo-be/modules/auth/domain/entities"
)

type OtpRepository interface {
	GetOtpRedis(ctx context.Context, username string) (data *entities.OTP, err error)
	SetOtpRedis(ctx context.Context, username string, data entities.OTP, exp time.Duration) (err error)
	IncrOtpAttemptRedis(ctx context.Context, username string) (res *int64, err error)
	GetOtpAttemptRedis(ctx context.Context, username string) (res *int64, err error)
}
