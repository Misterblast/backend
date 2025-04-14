package app

type AppError struct {
	Code      int
	ErrorCode string
	Message   string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string, errorCode ...string) *AppError {
	errCode := ""
	if len(errorCode) > 0 {
		errCode = errorCode[0]
	} else {
		errCode = "unknown-error"
	}
	return &AppError{
		Code:      code,
		Message:   message,
		ErrorCode: errCode,
	}
}

var (
	ErrBadRequest   = NewAppError(400, "bad-request", "bad request")
	ErrNotFound     = NewAppError(404, "resource-not-found", "resource not found")
	ErrInternal     = NewAppError(500, "unknown-error", "internal server error")
	ErrUnauthorized = NewAppError(401, "unauthorized", "unauthorized")
)
