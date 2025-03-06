package app

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrBadRequest   = NewAppError(400, "bad request")
	ErrNotFound     = NewAppError(404, "resource not found")
	ErrInternal     = NewAppError(500, "internal server error")
	ErrUnauthorized = NewAppError(401, "unauthorized")
)
