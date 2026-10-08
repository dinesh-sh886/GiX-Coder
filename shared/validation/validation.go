package validation

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator with custom validations.
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator with custom validations.
func NewValidator() (*Validator, error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.RegisterValidation("hostname", validateHostname); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("url", validateURL); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("duration", validateDuration); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("secret", validateSecret); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("tenant_id", validateTenantID); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("execution_id", validateExecutionID); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("correlation_id", validateCorrelationID); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("trace_id", validateTraceID); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("span_id", validateSpanID); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("capability_type", validateCapabilityType); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("resource_scope", validateResourceScope); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("jwt_token", validateJWTToken); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("base64", validateBase64); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("json", validateJSON); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("semver", validateSemVer); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("capability_grant", validateCapabilityGrant); err != nil {
		return nil, err
	}

	return &Validator{validator: v}, nil
}

// Validate validates a struct.
func (v *Validator) Validate(s any) error {
	return v.validator.Struct(s)
}

// ValidateVar validates a single variable.
func (v *Validator) ValidateVar(field any, tag string) error {
	return v.validator.Var(field, tag)
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Tag     string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors.
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// ValidateStruct validates a struct and returns detailed errors.
func (v *Validator) ValidateStruct(s any) ValidationErrors {
	err := v.validator.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors ValidationErrors
	validationErrors = append(validationErrors, ValidationError{
		Field:   "",
		Tag:     "validation",
		Value:   "",
		Message: err.Error(),
	})
	return validationErrors
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "must be a valid email address"
	case "url":
		return "must be a valid URL"
	case "hostname":
		return "must be a valid hostname"
	case "min":
		return fmt.Sprintf("must be at least %s", err.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", err.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters", err.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", err.Param())
	case "uuid":
		return "must be a valid UUID"
	case "duration":
		return "must be a valid duration (e.g., 30s, 5m, 1h)"
	case "secret":
		return "must be a valid secret reference"
	case "tenant_id":
		return "must be a valid tenant ID"
	case "execution_id":
		return "must be a valid execution ID"
	case "correlation_id":
		return "must be a valid correlation ID"
	case "trace_id":
		return "must be a valid trace ID"
	case "span_id":
		return "must be a valid span ID"
	case "capability_type":
		return "must be a valid capability type"
	case "resource_scope":
		return "must be a valid resource scope"
	case "jwt_token":
		return "must be a valid JWT token"
	case "base64":
		return "must be valid base64"
	case "json":
		return "must be valid JSON"
	case "semver":
		return "must be a valid semantic version"
	case "capability_grant":
		return "must be a valid capability grant"
	default:
		return fmt.Sprintf("validation failed: %s", err.Tag())
	}
}

// Custom validation functions

func validateHostname(fl validator.FieldLevel) bool {
	hostname := fl.Field().String()
	if len(hostname) > 255 {
		return false
	}
	if hostname == "" {
		return true
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return re.MatchString(hostname)
}

func validateURL(fl validator.FieldLevel) bool {
	urlStr := fl.Field().String()
	if urlStr == "" {
		return true
	}
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

func validateDuration(fl validator.FieldLevel) bool {
	duration := fl.Field().String()
	if duration == "" {
		return true
	}
	_, err := time.ParseDuration(duration)
	return err == nil
}

func validateSecret(fl validator.FieldLevel) bool {
	secret := fl.Field().String()
	if secret == "" {
		return true
	}
	return strings.HasPrefix(secret, "op://") || strings.HasPrefix(secret, "vault:")
}

func validateTenantID(fl validator.FieldLevel) bool {
	tenantID := fl.Field().String()
	if tenantID == "" {
		return true
	}
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]{1,64}$`, tenantID)
	if err != nil {
		return false
	}
	return matched
}

func validateExecutionID(fl validator.FieldLevel) bool {
	executionID := fl.Field().String()
	if executionID == "" {
		return true
	}
	matched, err := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, executionID)
	if err != nil {
		return false
	}
	return matched
}

func validateCorrelationID(fl validator.FieldLevel) bool {
	correlationID := fl.Field().String()
	if correlationID == "" {
		return true
	}
	matched, err := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, correlationID)
	if err != nil {
		return false
	}
	return matched
}

func validateTraceID(fl validator.FieldLevel) bool {
	traceID := fl.Field().String()
	if traceID == "" {
		return true
	}
	matched, err := regexp.MatchString(`^[0-9a-f]{32}$`, traceID)
	if err != nil {
		return false
	}
	return matched
}

func validateSpanID(fl validator.FieldLevel) bool {
	spanID := fl.Field().String()
	if spanID == "" {
		return true
	}
	matched, err := regexp.MatchString(`^[0-9a-f]{16}$`, spanID)
	if err != nil {
		return false
	}
	return matched
}

func validateCapabilityType(fl validator.FieldLevel) bool {
	capType := fl.Field().String()
	validTypes := map[string]bool{
		"filesystem.read":  true,
		"filesystem.write": true,
		"filesystem.list":  true,
		"shell.execute":    true,
		"git.read":         true,
		"git.write":        true,
		"test.execute":     true,
		"network.egress":   true,
		"mcp.invoke":       true,
	}
	return validTypes[capType]
}

func validateResourceScope(fl validator.FieldLevel) bool {
	scope := fl.Field().String()
	if scope == "" {
		return true
	}
	return len(scope) <= 512
}

func validateJWTToken(fl validator.FieldLevel) bool {
	token := fl.Field().String()
	if token == "" {
		return true
	}
	parts := strings.Split(token, ".")
	return len(parts) == 3
}

func validateBase64(fl validator.FieldLevel) bool {
	data := fl.Field().String()
	if data == "" {
		return true
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9+/]*={0,2}$`, data)
	if err != nil {
		return false
	}
	return matched
}

func validateJSON(fl validator.FieldLevel) bool {
	jsonStr := fl.Field().String()
	if jsonStr == "" {
		return true
	}
	return strings.HasPrefix(strings.TrimSpace(jsonStr), "{") || strings.HasPrefix(strings.TrimSpace(jsonStr), "[")
}

func validateSemVer(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	if version == "" {
		return true
	}
	semverRegex := regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	return semverRegex.MatchString(version)
}

func validateCapabilityGrant(fl validator.FieldLevel) bool {
	return !fl.Field().IsNil()
}

// SanitizeInput sanitizes input to prevent injection attacks.
func SanitizeInput(input string) string {
	input = strings.ReplaceAll(input, "\x00", "")
	var result strings.Builder
	for _, r := range input {
		if r >= 32 || r == '\n' || r == '\t' || r == '\r' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// SanitizePath sanitizes a file path to prevent traversal.
func SanitizePath(path string) string {
	path = strings.ReplaceAll(path, "\x00", "")
	path = filepath.Clean(path)
	if strings.HasPrefix(path, "..") || strings.Contains(path, "/..") {
		return ""
	}
	return path
}
