package config

type Otp struct {
	Expiration    int64 `yaml:"expiration"`
	MaxAttempt    int64 `yaml:"maxAttempts"`
	RetryInterval int64 `yaml:"retryInterval"`
}
