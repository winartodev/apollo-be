package response

import (
	domainError "github.com/winartodev/apollo-be/internal/domain/error"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

const (
	ErrInvalidRequestPayload = "invalid request payload %v"
)

// FieldError represents validation error for a specific field
// swagger:model
type FieldError struct {
	// The field name with error
	Field string `json:"field"`

	// Error message for the field
	Message string `json:"message"`
}

// Response represents a standard API response format
// swagger:model
type Response struct {
	// Indicates if the request was successful
	Success bool `json:"success"`

	// Optional message describing the result
	Message string `json:"message,omitempty"`

	// The main data payload
	Data interface{} `json:"data,omitempty"`

	// Additional metadata about the response
	Meta interface{} `json:"meta,omitempty"`
}

// ErrorResponse represents a standard API error response format
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Indicates if the request was successful
	Success bool `json:"success"`

	// Optional message describing the result
	Message string `json:"message,omitempty"`

	// General error message when no field-specific errors exist
	// Example: "Invalid credentials"
	Error interface{} `json:"error,omitempty"`
}

// ValidationErrorResponse represents validation errors with field-specific details
// swagger:model ValidationErrorResponse
type ValidationErrorResponse struct {
	// Indicates if the request was successful
	Success bool `json:"success"`

	// Optional message describing the result
	Message string `json:"message,omitempty"`

	// Field-specific validation errors
	// Required: true
	Error []FieldError `json:"error,omitempty"`
}

// PaginateResponse represents a paginated API response
// swagger:model
type PaginateResponse struct {
	// The paginated data
	Data interface{} `json:"data"`

	// Total number of items available
	Total int64 `json:"total"`

	// Current page number
	Page int `json:"page"`

	// Number of items per page
	PerPage int `json:"per_page"`

	// Total number of pages
	TotalPages int `json:"total_pages"`
}

func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}, meta interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func FailedResponse(c echo.Context, statusCode int, err error) error {
	statusCode = domainError.GetHTTPStatusFromError(err)
	return c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   err.Error(),
	})
}

func ValidationErrResponse(c echo.Context, err error) error {
	var FieldErrors []FieldError
	for _, e := range err.(validator.ValidationErrors) {
		FieldErrors = append(FieldErrors, FieldError{
			Field:   e.Field(),
			Message: getValidationErrorMessage(e),
		})
	}

	return c.JSON(http.StatusUnprocessableEntity, ValidationErrorResponse{
		Success: false,
		Error:   FieldErrors,
	})
}

func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "This field must be at least " + fe.Param() + " characters"
	case "max":
		return "This field must be at most " + fe.Param() + " characters"
	case "len":
		return "This field must be exactly " + fe.Param() + " characters"
	case "numeric":
		return "This field must be numeric"
	case "alpha":
		return "This field must contain only letters"
	case "alphanum":
		return "This field must contain only letters and numbers"
	default:
		return "Invalid value"
	}
}
