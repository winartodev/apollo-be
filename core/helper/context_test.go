package helper

import (
	"context"
	"errors"
	"testing"
)

func TestSaveUserIDToContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		userID   int64
		expected int64
	}{
		{
			name:     "Save valid user ID to background context",
			ctx:      context.Background(),
			userID:   123456,
			expected: 123456,
		},
		{
			name:     "Save zero user ID",
			ctx:      context.Background(),
			userID:   0,
			expected: 0,
		},
		{
			name:     "Save negative user ID",
			ctx:      context.Background(),
			userID:   -1,
			expected: -1,
		},
		{
			name:     "Save large user ID",
			ctx:      context.Background(),
			userID:   9223372036854775807, // max int64
			expected: 9223372036854775807,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := SaveUserIDToContext(&tt.ctx, tt.userID)

			if result == nil {
				t.Fatal("SaveUserIDToContext returned nil context")
			}

			// Verify the value was saved correctly
			savedValue := result.Value(userIdKey)
			if savedValue == nil {
				t.Fatal("User ID not found in returned context")
			}

			savedUserID, ok := savedValue.(int64)
			if !ok {
				t.Fatal("Saved value is not int64")
			}

			if savedUserID != tt.expected {
				t.Errorf("Expected user ID %d, got %d", tt.expected, savedUserID)
			}
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name          string
		setupContext  func() context.Context
		expectedID    int64
		expectedError error
		shouldError   bool
	}{
		{
			name: "Get valid user ID from context",
			setupContext: func() context.Context {
				ctx := context.Background()
				return SaveUserIDToContext(&ctx, 123456)
			},
			expectedID:    123456,
			expectedError: nil,
			shouldError:   false,
		},
		{
			name: "Get zero user ID from context",
			setupContext: func() context.Context {
				ctx := context.Background()
				return SaveUserIDToContext(&ctx, 0)
			},
			expectedID:    0,
			expectedError: nil,
			shouldError:   false,
		},
		{
			name: "Get user ID from empty context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedID:    0,
			expectedError: errUserIDNotFound,
			shouldError:   true,
		},
		{
			name: "Get user ID when context has wrong data type",
			setupContext: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, userIdKey, "not-an-int64")
			},
			expectedID:    0,
			expectedError: errInvalidUserIDDataType,
			shouldError:   true,
		},
		{
			name: "Get user ID when context has nil value",
			setupContext: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, userIdKey, nil)
			},
			expectedID:    0,
			expectedError: errUserIDNotFound,
			shouldError:   true,
		},
		{
			name: "Get user ID when context has different numeric type",
			setupContext: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, userIdKey, int32(12345))
			},
			expectedID:    0,
			expectedError: errInvalidUserIDDataType,
			shouldError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()

			userID, err := GetUserIDFromContext(ctx)

			if tt.shouldError {
				if err == nil {
					t.Fatal("Expected an error but got none")
				}
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
				if userID != tt.expectedID {
					t.Errorf("Expected user ID %d, got %d", tt.expectedID, userID)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}
				if userID != tt.expectedID {
					t.Errorf("Expected user ID %d, got %d", tt.expectedID, userID)
				}
			}
		})
	}
}

func TestSaveAndGetUserIDIntegration(t *testing.T) {
	testCases := []int64{
		1,
		123456789,
		0,
		-1,
		9223372036854775807,  // max int64
		-9223372036854775808, // min int64
	}

	for _, userID := range testCases {
		t.Run("Integration test with user ID", func(t *testing.T) {
			ctx := context.Background()

			contextWithUserID := SaveUserIDToContext(&ctx, userID)
			retrievedUserID, err := GetUserIDFromContext(contextWithUserID)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if retrievedUserID != userID {
				t.Errorf("Expected user ID %d, got %d", userID, retrievedUserID)
			}
		})
	}
}
