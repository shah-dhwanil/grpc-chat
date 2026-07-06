package pkgerror

import "fmt"

type ErrorType string
const (
	Unauthorized ErrorType = "UNAUTHORIZED"
	ResourceAlreadyExists ErrorType = "RESOURCE_ALREADY_EXISTS"
	ResourceNotFound      ErrorType = "RESOURCE_NOT_FOUND"
	Validation ErrorType = "VALIDATION_ERROR"
	Internal   ErrorType = "INTERNAL_ERROR"
	Unknown    ErrorType = "UNKNOWN_ERROR"
)

// AppError is a custom error type representing any error occured during the execution of the application.
type AppError struct {
	Type    ErrorType  // Type of the error, e.g., validation_error, unknown_error, etc.
	Name   string     // A short, human-readable summary of the error.
	Detail string     // A detailed description of the error, providing more context and information about what went wrong.
	Context map[string]any  // Additional context about the error, if applicable.
	wrappedError error // The original error that caused this AppError, if applicable. This field is not included in JSON serialization.
}

func (e *AppError) Error() string {
    if e.wrappedError != nil {
        return fmt.Sprintf("%s [%s]: %s - %s", e.Type, e.Name, e.Detail, e.wrappedError.Error())
    }
    return fmt.Sprintf("%s [%s]: %s", e.Type, e.Name, e.Detail)
}


func (e *AppError) Unwrap() error {
	return e.wrappedError
}

func NewAppError(errType ErrorType, name string, message string, context map[string]any, wrappedError error) *AppError {
	return &AppError{
		Type:         errType,
		Name:        name,
		Detail:      message,
		Context:      context,
		wrappedError: wrappedError,
	}
}