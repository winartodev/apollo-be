package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/modules/auth/domain/entities"
	"github.com/winartodev/apollo-be/modules/auth/domain/repository"
)

const (
	otpRedisKey         = "otp:%s"
	otpAttemptsRedisKey = "otp_attempts:%s"
)

type OtpRepositoryImpl struct {
	*helper.RedisUtil
}

func NewOtpRepository(redis *helper.RedisUtil) (repository.OtpRepository, error) {
	return &OtpRepositoryImpl{
		RedisUtil: redis,
	}, nil
}

func (r *OtpRepositoryImpl) SetOtpRedis(ctx context.Context, username string, data entities.OTP, exp time.Duration) (err error) {
	key := fmt.Sprintf(otpRedisKey, username)
	return r.RedisUtil.SetEx(ctx, key, data, exp)
}

func (r *OtpRepositoryImpl) IncrOtpAttemptRedis(ctx context.Context, username string) (res *int64, err error) {
	key := fmt.Sprintf(otpAttemptsRedisKey, username)
	val, err := r.RedisUtil.IncrBy(ctx, key, 1)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (r *OtpRepositoryImpl) GetOtpRedis(ctx context.Context, username string) (data *entities.OTP, err error) {
	key := fmt.Sprintf(otpRedisKey, username)
	err = r.RedisUtil.Get(ctx, key, &data)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	return data, nil
}

func (r *OtpRepositoryImpl) GetOtpAttemptRedis(ctx context.Context, username string) (res *int64, err error) {
	key := fmt.Sprintf(otpAttemptsRedisKey, username)
	err = r.RedisUtil.Get(ctx, key, &res)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}
