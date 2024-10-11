package errors

type ForbiddenError struct {
	err error
}

func (*ForbiddenError) StatusCode() int {
	return 403
}

func (e *ForbiddenError) ErrorCode() string {
	return "FORBIDDEN"
}

func (e *ForbiddenError) Error() string {
	return e.err.Error()
}

func NewForbiddenError(err error) *ForbiddenError {
	return &ForbiddenError{err: err}
}

type Unauthorized struct{}

func (*Unauthorized) StatusCode() int {
	return 401
}

func (e *Unauthorized) ErrorCode() string {
	return "UNAUTHORIZED"
}

func (e *Unauthorized) Error() string {
	return "Не аутентифицирован"
}
