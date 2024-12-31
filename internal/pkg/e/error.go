package e

type ApiError interface {
	Error() string
	Code() int
}

type CustomApiError struct {
	ErrCode    int  
	ErrMessage string
}

func NewApiError(code int, message string) *CustomApiError {
	return &CustomApiError{
		ErrCode:    code,
		ErrMessage: message,
	}
}

func (e *CustomApiError) Error() string {
	return e.ErrMessage
}

func (e *CustomApiError) Code() int {
	return e.ErrCode
}