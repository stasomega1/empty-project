package model

type MyError struct {
	UserMessage string
	HttpStatus  int
	DebugInfo   string
}

func NewError(errText string, httpStatus int, debugInfo string) *MyError {
	return &MyError{
		UserMessage: errText,
		HttpStatus:  httpStatus,
		DebugInfo:   debugInfo,
	}
}

func (e *MyError) Error() string {
	return e.UserMessage
}
