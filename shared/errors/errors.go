package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents a stable, machine-readable error code.
type ErrorCode string

const (
	// Generic errors
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeRateLimited  ErrorCode = "RATE_LIMITED"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodePrecondition ErrorCode = "PRECONDITION_FAILED"
	ErrCodeTimeout      ErrorCode = "TIMEOUT"
	ErrCodeUnavailable  ErrorCode = "SERVICE_UNAVAILABLE"

	// Authentication/Authorization
	ErrCodeAuthFailed       ErrorCode = "AUTH_FAILED"
	ErrCodeAuthExpired      ErrorCode = "AUTH_EXPIRED"
	ErrCodeAuthInvalid      ErrorCode = "AUTH_INVALID"
	ErrCodeAuthzDenied      ErrorCode = "AUTHZ_DENIED"
	ErrCodeCapabilityDenied ErrorCode = "CAPABILITY_DENIED"

	// Resource/Execution
	ErrCodeResourceExceeded ErrorCode = "RESOURCE_EXCEEDED"
	ErrCodeSandboxError     ErrorCode = "SANDBOX_ERROR"
	ErrCodeBudgetExceeded   ErrorCode = "BUDGET_EXCEEDED"
	ErrCodeExecutionFailed  ErrorCode = "EXECUTION_FAILED"
	ErrCodeCheckpointFailed ErrorCode = "CHECKPOINT_FAILED"

	// External dependencies
	ErrCodeDependencyError ErrorCode = "DEPENDENCY_ERROR"
	ErrCodeProviderError   ErrorCode = "PROVIDER_ERROR"
)

// HTTPStatus returns the appropriate HTTP status code for the error code.
func (c ErrorCode) HTTPStatus() int {
	switch c {
	case ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeAuthFailed, ErrCodeAuthExpired, ErrCodeAuthInvalid:
		return http.StatusUnauthorized
	case ErrCodeForbidden, ErrCodeAuthzDenied, ErrCodeCapabilityDenied:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodePrecondition:
		return http.StatusPreconditionFailed
	case ErrCodeRateLimited:
		return http.StatusTooManyRequests
	case ErrCodeTimeout:
		return http.StatusGatewayTimeout
	case ErrCodeUnavailable, ErrCodeDependencyError, ErrCodeProviderError:
		return http.StatusServiceUnavailable
	case ErrCodeResourceExceeded, ErrCodeBudgetExceeded:
		return http.StatusInsufficientStorage
	case ErrCodeInternal, ErrCodeSandboxError, ErrCodeExecutionFailed, ErrCodeCheckpointFailed:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// AppError is the standard application error type.
type AppError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
	Cause   error
	Stack   string
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithDetail adds a detail to the error.
func (e *AppError) WithDetail(key string, value any) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

// WithCause wraps an error with additional context.
func WithCause(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return &AppError{Code: code, Message: message}
	}
	var appErr *AppError
	if As(err, &appErr) {
		return appErr.WithDetail("wrapped", message)
	}
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// New creates a new AppError.
func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Is checks if an error matches a specific error code.
func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// As finds the first error in the chain that matches the target type.
func As(err error, target any) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*AppError); ok {
		switch t := target.(type) {
		case **AppError:
			*t = e
			return true
		case *ErrorCode:
			*t = e.Code
			return true
		}
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		switch t := target.(type) {
		case **AppError:
			*t = appErr
			return true
		case *ErrorCode:
			*t = appErr.Code
			return true
		}
		return true
	}
	return false
}

// Sentinel errors for common control flow.
var (
	ErrNotFound     = New(ErrCodeNotFound, "resource not found")
	ErrUnauthorized = New(ErrCodeUnauthorized, "unauthorized")
	ErrForbidden    = New(ErrCodeForbidden, "forbidden")
	ErrValidation   = New(ErrCodeValidation, "validation failed")
	ErrRateLimited  = New(ErrCodeRateLimited, "rate limited")
	ErrTimeout      = New(ErrCodeTimeout, "operation timed out")
	ErrInternal     = New(ErrCodeInternal, "internal server error")
	ErrUnavailable  = New(ErrCodeUnavailable, "service unavailable")
)

// GRPCCode maps ErrorCode to gRPC status codes.
func (c ErrorCode) GRPCCode() int {
	switch c {
	case ErrCodeValidation:
		return 3 // InvalidArgument
	case ErrCodeUnauthorized, ErrCodeAuthFailed, ErrCodeAuthExpired, ErrCodeAuthInvalid:
		return 16 // Unauthenticated
	case ErrCodeForbidden, ErrCodeAuthzDenied, ErrCodeCapabilityDenied:
		return 7 // PermissionDenied
	case ErrCodeNotFound:
		return 5 // NotFound
	case ErrCodeConflict:
		return 6 // AlreadyExists
	case ErrCodePrecondition:
		return 9 // FailedPrecondition
	case ErrCodeRateLimited:
		return 8 // ResourceExhausted
	case ErrCodeTimeout:
		return 4 // DeadlineExceeded
	case ErrCodeUnavailable, ErrCodeDependencyError, ErrCodeProviderError:
		return 14 // Unavailable
	case ErrCodeInternal, ErrCodeSandboxError, ErrCodeExecutionFailed, ErrCodeCheckpointFailed:
		return 2 // Unknown
	default:
		return 2 // Unknown
	}
}
