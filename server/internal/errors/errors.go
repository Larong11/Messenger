package errors

type ErrorCode int

/*
создан для того, чтобы когда ошибка возвращается из use_cases сразу понимать
какого рода эта ошибка, и без if возвращать нужный http ответ
*/

const (
	CodeBadRequest ErrorCode = iota + 1
	CodeInternal
)

type AppError struct {
	Code    ErrorCode
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequest(msg string) *AppError {
	return &AppError{Code: CodeBadRequest, Message: msg}
}

func NewInternal(msg string) *AppError {
	return &AppError{Code: CodeInternal, Message: msg}
}
