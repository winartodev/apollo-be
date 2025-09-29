package context

import (
	"context"
	"errors"
)

type ContextKey string

const (
	UserIdKey      ContextKey = "user_id"
	AppPlatformKey ContextKey = "application_platform"
)

var (
	errUserIDNotFound        = errors.New("user ID not found in context")
	errInvalidUserIDDataType = errors.New("user ID is not of type int64")

	errAppPlatformNotFound = errors.New("app platform not found in context")
	errInvalidAppPlatform  = errors.New("invalid app platform")
)

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	value := ctx.Value(UserIdKey)
	if value == nil {
		return 0, errUserIDNotFound
	}

	userID, ok := value.(int64)
	if !ok {
		return 0, errInvalidUserIDDataType
	}

	return userID, nil
}

func GetAppPlatformFromContext(ctx context.Context) (string, error) {
	value := ctx.Value(AppPlatformKey)
	if value == nil {
		return "", errAppPlatformNotFound
	}

	appPlatform, ok := value.(string)
	if !ok {
		return "", errInvalidAppPlatform
	}

	return appPlatform, nil
}
