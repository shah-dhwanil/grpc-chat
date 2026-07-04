package pkgerror

func NewUnknownError(wrappedError error, title, message string, context map[string]any) *AppError {
	return &AppError{
		Type:        Unknown,
		Name:       title,
		Detail:     message,
		Context:     context,
		wrappedError: wrappedError,
	}
}

func NewInternalError(wrappedError error, title, message string, context map[string]any) *AppError {
	return &AppError{
		Type:        Internal,
		Name:       title,
		Detail:     message,
		Context:     context,
		wrappedError: wrappedError,
	}
}