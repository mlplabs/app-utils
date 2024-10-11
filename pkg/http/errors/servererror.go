package errors

type ServerError struct {
	Err error
}

func (e *ServerError) ErrorCode() string {
	return "SERVER_UNEXPECTED"
}

func (*ServerError) StatusCode() int {
	return 500
}

func (e *ServerError) Error() string {
	return "Ошибка сервиса"
}

func NewServerError(err error) *ServerError {
	return &ServerError{Err: err}
}
