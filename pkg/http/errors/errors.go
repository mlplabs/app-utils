package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type CommonError interface {
	StatusCode() int
	ErrorCode() string
	Error() string
}

type ResponseError struct {
	Error Response `json:"error"`
}
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Service string      `json:"service,omitempty"` // Need for understanding to get name service with error
}

func parseError(err error) CommonError {
	var e CommonError
	if !errors.As(err, &e) {
		e = NewServerError(err)
	}
	return e
}

func SetError(w http.ResponseWriter, err error) {
	commonErr := parseError(err)
	w.WriteHeader(commonErr.StatusCode()) // менять от ошибки
	data := ResponseError{
		Error: Response{
			Code:    commonErr.ErrorCode(),
			Message: commonErr.Error(),
			Data:    nil,
			Service: "",
		},
	}
	body, err := json.Marshal(data)
	if err != nil {
		SetError(w, err)
	}
	w.Write(body)
}
