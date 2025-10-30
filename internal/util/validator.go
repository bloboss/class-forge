package util

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Error implements the error interface
func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (ves ValidationErrors) Error() string {
	if len(ves) == 0 {
		return ""
	}
	messages := make([]string, len(ves))
	for i, ve := range ves {
		messages[i] = ve.Error()
	}
	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are validation errors
func (ves ValidationErrors) HasErrors() bool {
	return len(ves) > 0
}

// Validator provides validation utilities
type Validator struct {
	errors ValidationErrors
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: make(ValidationErrors, 0),
	}
}

// AddError adds a validation error
func (v *Validator) AddError(field, message, code string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns all validation errors
func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

// ValidateRequired validates that a field is not empty
func (v *Validator) ValidateRequired(field, value, displayName string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, fmt.Sprintf("%s is required", displayName), "VALIDATION_MISSING_REQUIRED_FIELD")
	}
}

// ValidateLength validates string length
func (v *Validator) ValidateLength(field, value, displayName string, min, max int) {
	length := utf8.RuneCountInString(value)
	if length < min {
		v.AddError(field, fmt.Sprintf("%s must be at least %d characters", displayName, min), "VALIDATION_TOO_SHORT")
	}
	if max > 0 && length > max {
		v.AddError(field, fmt.Sprintf("%s must be no more than %d characters", displayName, max), "VALIDATION_TOO_LONG")
	}
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(field, email, displayName string) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if email != "" && !emailRegex.MatchString(email) {
		v.AddError(field, fmt.Sprintf("%s must be a valid email address", displayName), "VALIDATION_INVALID_FORMAT")
	}
}

// ValidateSlug validates slug format (lowercase, alphanumeric, hyphens)
func (v *Validator) ValidateSlug(field, slug, displayName string) {
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	if slug != "" && !slugRegex.MatchString(slug) {
		v.AddError(field, fmt.Sprintf("%s must be a valid slug (lowercase letters, numbers, and hyphens only)", displayName), "VALIDATION_INVALID_FORMAT")
	}
}

// ValidateURL validates URL format
func (v *Validator) ValidateURL(field, url, displayName string) {
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+(?:[/?#].*)?$`)
	if url != "" && !urlRegex.MatchString(url) {
		v.AddError(field, fmt.Sprintf("%s must be a valid URL", displayName), "VALIDATION_INVALID_FORMAT")
	}
}

// ValidateDateTime validates RFC3339 datetime format
func (v *Validator) ValidateDateTime(field, datetime, displayName string) {
	if datetime != "" {
		if _, err := time.Parse(time.RFC3339, datetime); err != nil {
			v.AddError(field, fmt.Sprintf("%s must be a valid RFC3339 datetime", displayName), "VALIDATION_INVALID_DATE")
		}
	}
}

// ValidateFutureDate validates that a datetime is in the future
func (v *Validator) ValidateFutureDate(field, datetime, displayName string) {
	if datetime != "" {
		if t, err := time.Parse(time.RFC3339, datetime); err == nil {
			if !t.After(time.Now()) {
				v.AddError(field, fmt.Sprintf("%s must be in the future", displayName), "VALIDATION_INVALID_DATE")
			}
		}
	}
}

// ValidateEnum validates that a value is in a set of allowed values
func (v *Validator) ValidateEnum(field, value, displayName string, allowedValues []string) {
	if value != "" {
		allowed := false
		for _, allowed_value := range allowedValues {
			if value == allowed_value {
				allowed = true
				break
			}
		}
		if !allowed {
			v.AddError(field, fmt.Sprintf("%s must be one of: %s", displayName, strings.Join(allowedValues, ", ")), "VALIDATION_INVALID_INPUT")
		}
	}
}

// ValidatePositiveInt validates that an integer is positive
func (v *Validator) ValidatePositiveInt(field string, value int, displayName string) {
	if value <= 0 {
		v.AddError(field, fmt.Sprintf("%s must be a positive number", displayName), "VALIDATION_INVALID_INPUT")
	}
}

// ValidateRange validates that an integer is within a range
func (v *Validator) ValidateRange(field string, value, min, max int, displayName string) {
	if value < min || value > max {
		v.AddError(field, fmt.Sprintf("%s must be between %d and %d", displayName, min, max), "VALIDATION_INVALID_INPUT")
	}
}