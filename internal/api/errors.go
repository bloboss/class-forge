package api

// Error code taxonomy as defined in design.md Section 6.2
const (
	// Authentication Errors (AUTH_*)
	ErrAuthMissingToken = "AUTH_MISSING_TOKEN"
	ErrAuthInvalidToken = "AUTH_INVALID_TOKEN"
	ErrAuthExpiredToken = "AUTH_EXPIRED_TOKEN"

	// Authorization Errors (AUTHZ_*)
	ErrAuthzForbidden               = "AUTHZ_FORBIDDEN"
	ErrAuthzInsufficientPermissions = "AUTHZ_INSUFFICIENT_PERMISSIONS"

	// Validation Errors (VALIDATION_*)
	ErrValidationInvalidInput  = "VALIDATION_INVALID_INPUT"
	ErrValidationMissingField  = "VALIDATION_MISSING_REQUIRED_FIELD"
	ErrValidationInvalidFormat = "VALIDATION_INVALID_FORMAT"
	ErrValidationInvalidDate   = "VALIDATION_INVALID_DATE"
	ErrValidationTooShort      = "VALIDATION_TOO_SHORT"
	ErrValidationTooLong       = "VALIDATION_TOO_LONG"

	// Resource Errors (RESOURCE_*)
	ErrResourceNotFound      = "RESOURCE_NOT_FOUND"
	ErrResourceConflict      = "RESOURCE_CONFLICT"
	ErrResourceAlreadyExists = "RESOURCE_ALREADY_EXISTS"

	// Business Logic Errors (BUSINESS_*)
	ErrBusinessDeadlinePassed   = "BUSINESS_DEADLINE_PASSED"
	ErrBusinessAlreadyAccepted  = "BUSINESS_ALREADY_ACCEPTED"
	ErrBusinessRosterNotFound   = "BUSINESS_ROSTER_NOT_FOUND"
	ErrBusinessTeamSizeExceeded = "BUSINESS_TEAM_SIZE_EXCEEDED"
	ErrBusinessTemplateNotFound = "BUSINESS_TEMPLATE_NOT_FOUND"

	// Integration Errors (INTEGRATION_*)
	ErrIntegrationForgejoAPI         = "INTEGRATION_FORGEJO_API_ERROR"
	ErrIntegrationForgejoRateLimited = "INTEGRATION_FORGEJO_RATE_LIMITED"
	ErrIntegrationForgejoUnavailable = "INTEGRATION_FORGEJO_UNAVAILABLE"
	ErrIntegrationDatabase           = "INTEGRATION_DATABASE_ERROR"

	// System Errors (SYSTEM_*)
	ErrSystemInternal    = "SYSTEM_INTERNAL_ERROR"
	ErrSystemUnavailable = "SYSTEM_UNAVAILABLE"
	ErrSystemTimeout     = "SYSTEM_TIMEOUT"
)

// ErrorMessages provides human-readable messages for error codes
var ErrorMessages = map[string]string{
	// Authentication Errors
	ErrAuthMissingToken: "Authorization token is required",
	ErrAuthInvalidToken: "Invalid authorization token",
	ErrAuthExpiredToken: "Authorization token has expired",

	// Authorization Errors
	ErrAuthzForbidden:               "Access forbidden",
	ErrAuthzInsufficientPermissions: "Insufficient permissions for this operation",

	// Validation Errors
	ErrValidationInvalidInput:  "Invalid input provided",
	ErrValidationMissingField:  "Required field is missing",
	ErrValidationInvalidFormat: "Invalid format",
	ErrValidationInvalidDate:   "Invalid date format or value",
	ErrValidationTooShort:      "Value is too short",
	ErrValidationTooLong:       "Value is too long",

	// Resource Errors
	ErrResourceNotFound:      "Requested resource not found",
	ErrResourceConflict:      "Resource conflict detected",
	ErrResourceAlreadyExists: "Resource already exists",

	// Business Logic Errors
	ErrBusinessDeadlinePassed:   "Assignment deadline has passed",
	ErrBusinessAlreadyAccepted:  "Assignment has already been accepted",
	ErrBusinessRosterNotFound:   "Student not found in classroom roster",
	ErrBusinessTeamSizeExceeded: "Team size limit exceeded",
	ErrBusinessTemplateNotFound: "Assignment template repository not found",

	// Integration Errors
	ErrIntegrationForgejoAPI:         "Forgejo API error",
	ErrIntegrationForgejoRateLimited: "Forgejo API rate limit exceeded",
	ErrIntegrationForgejoUnavailable: "Forgejo service unavailable",
	ErrIntegrationDatabase:           "Database operation failed",

	// System Errors
	ErrSystemInternal:    "Internal server error",
	ErrSystemUnavailable: "Service temporarily unavailable",
	ErrSystemTimeout:     "Request timeout",
}

// GetErrorMessage returns the human-readable message for an error code
func GetErrorMessage(code string) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "Unknown error"
}
