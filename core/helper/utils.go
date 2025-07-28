package helper

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

const (
	errorReadYamlFile  = "error while reading yaml file %v"
	errorUnmarshalFile = "error while unmarshalling file %v"
	errInvalidPath     = "invalid path"
)

func ReadYaml(filename string, data interface{}) error {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return errors.New(fmt.Sprintf(errorReadYamlFile, err))
	}

	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		return errors.New(fmt.Sprintf(errorUnmarshalFile, err))
	}

	return nil
}

func ternary(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}

	return falseVal
}

func GetCompletePath(path string) (completePath string, err error) {
	completePath, err = filepath.Abs(path)
	if err != nil {
		return "", errors.New(errInvalidPath)
	}

	return completePath, nil
}

func CurrentOS(os string) bool {
	return runtime.GOOS == os
}

func BuildRedirectionLink[T any](
	ctx context.Context,
	operation T,
	mobileRedirect func(T) string,
	webRedirect func(T) string,
) string {
	platform, err := getAppPlatformFromContext(ctx)
	if err != nil {
		return webRedirect(operation)
	}

	switch platform {
	case "android", "ios":
		return mobileRedirect(operation)
	case "web":
		return webRedirect(operation)
	default:
		return webRedirect(operation)
	}
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetFirstElement(codes []string) string {
	if len(codes) > 0 {
		return codes[0]
	}
	return ""
}
