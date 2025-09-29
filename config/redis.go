package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
	PoolSize int    `yaml:"poolSize"`
}

func (rc *RedisConfig) SetupConnection() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rc.Host, rc.Port),
		Password: rc.Password,
		DB:       rc.Database,
		PoolSize: rc.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Health check with retry logic
	var err error
	for i := 0; i < 3; i++ {
		if err = client.Ping(ctx).Err(); err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis after 3 attempts: %v", err)
	}

	return client, nil
}
