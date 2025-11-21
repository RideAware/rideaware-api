package errors

type AppError struct {
	Code    int
	Message string
	Details string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

var (
	ErrUnauthorized = NewAppError(401, "Unauthorized", "")
	ErrNotFound     = NewAppError(404, "Not Found", "")
	ErrBadRequest   = NewAppError(400, "Bad Request", "")
	ErrInternal     = NewAppError(500, "Internal Server Error", "")
)