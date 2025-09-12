package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/winartodev/apollo-be/helper"
)

const (
	developmentConfigPath = "files/apollo.development.yaml"

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

	Jwt Jwt `yaml:"jwt"`

	Redis RedisConfig `yaml:"redis"`

	SMTP SMTPConfig `yaml:"smtp"`

	OTP Otp `yaml:"otp"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := helper.ReadYaml(developmentConfigPath, &cfg)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(errorLoadConfig, err))
	}

	if err = os.Setenv(JwtAccessTokenSecretKey, cfg.Jwt.AccessTokenSecret); err != nil {
		return nil, err
	}

	if err = os.Setenv(JwtRefreshTokenSecretKey, cfg.Jwt.RefreshTokenSecret); err != nil {
		return nil, err
	}

	return &cfg, nil
}
