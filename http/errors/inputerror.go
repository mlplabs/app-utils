package errors

type InvalidInputData struct {
	err error
}

func (*InvalidInputData) StatusCode() int {
	return 400
}

func (*InvalidInputData) ErrorCode() string {
	return "INVALID_INPUT_DATA"
}

func (e *InvalidInputData) ErrorText() string {
	return e.err.Error()
}

func NewInvalidInputData(err error) *InvalidInputData {
	return &InvalidInputData{err: err}
}

type BadRequest struct {
	err error
}

func (*BadRequest) StatusCode() int {
	return 400
}

func (*BadRequest) ErrorCode() string {
	return "BAD_REQUEST"
}

func (e *BadRequest) ErrorText() string {
	return e.err.Error()
}

func NewBadRequest(err error) *BadRequest {
	return &BadRequest{err: err}
}
