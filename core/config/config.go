package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/winartodev/apollo-be/core/helper"
)

const (
	developmentConfigPath = "core/files/apollo.development.yaml"

	errorLoadConfig = "error while loading config file %v"
)

type Config struct {
	App struct {
		Name string `yaml:"name"`
	} `yaml:"app"`

	Http struct {
		Port string `yaml:"port"`
	}

	Database Database `yaml:"database"`

	Jwt struct {
		AccessTokenSecret  string `yaml:"accessTokenSecret"`
		RefreshTokenSecret string `yaml:"refreshTokenSecret"`
	} `yaml:"jwt"`

	Redis RedisConfig `yaml:"redis"`

	Smtp helper.SmtpConfig `yaml:"smtp"`

	OTP helper.OtpConfig `yaml:"otp"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := helper.ReadYaml(developmentConfigPath, &cfg)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(errorLoadConfig, err))
	}

	if err = os.Setenv(helper.JwtAccessTokenSecretKey, cfg.Jwt.AccessTokenSecret); err != nil {
		return nil, err
	}

	if err = os.Setenv(helper.JwtRefreshTokenSecretKey, cfg.Jwt.RefreshTokenSecret); err != nil {
		return nil, err
	}

	return &cfg, nil
}
