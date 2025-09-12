package helper

type OtpConfig struct {
	Expiration    int64 `yaml:"expiration"`
	MaxAttempt    int64 `yaml:"maxAttempt"`
	RetryInterval int64 `yaml:"retryInterval"`
}
