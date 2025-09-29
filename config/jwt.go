package config

const (
	JwtAccessTokenSecretKey  = "JWT_ACCESS_TOKEN_SECRET_KEY"
	JwtRefreshTokenSecretKey = "JWT_REFRESH_TOKEN_SECRET_KEY"
)

type Jwt struct {
	AccessTokenSecret  string `yaml:"accessTokenSecret"`
	RefreshTokenSecret string `yaml:"refreshTokenSecret"`
}
