package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	redisInfra "github.com/winartodev/apollo-be/infrastructure/redis"
	"github.com/winartodev/apollo-be/modules/auth/domain/entities"
	"github.com/winartodev/apollo-be/modules/auth/domain/repository"
)

const (
	otpRedisKey         = "otp:%s"
	otpAttemptsRedisKey = "otp_attempts:%s"
)

type OtpRepositoryImpl struct {
	*redisInfra.Redis
}

func NewOtpRepository(redisClient *redisInfra.Redis) (repository.OtpRepository, error) {
	return &OtpRepositoryImpl{
		Redis: redisClient,
	}, nil
}

func (r *OtpRepositoryImpl) SetOtpRedis(ctx context.Context, username string, data entities.OTP, exp time.Duration) (err error) {
	key := fmt.Sprintf(otpRedisKey, username)
	return r.Redis.SetEx(ctx, key, data, exp)
}

func (r *OtpRepositoryImpl) IncrOtpAttemptRedis(ctx context.Context, username string) (res *int64, err error) {
	key := fmt.Sprintf(otpAttemptsRedisKey, username)

	val, err := r.Redis.IncrBy(ctx, key, 1)
	if err != nil {
		return nil, err
	}

	err = r.Redis.Expire(ctx, key, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (r *OtpRepositoryImpl) GetOtpRedis(ctx context.Context, username string) (data *entities.OTP, err error) {
	key := fmt.Sprintf(otpRedisKey, username)
	err = r.Redis.Get(ctx, key, &data)
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
	err = r.Redis.Get(ctx, key, &res)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}
