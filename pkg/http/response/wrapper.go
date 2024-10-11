package response

import (
	"encoding/json"
	"github.com/mlplabs/app-utils/pkg/http/errors"
	"net/http"
)

type PlainData struct {
	Data interface{} `json:"data"`
}

type Wrapper struct{}

func NewWrapper() *Wrapper {
	return &Wrapper{}
}

func (rw *Wrapper) response(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	if data != nil {
		wrapData := PlainData{
			Data: data,
		}
		body, err := json.Marshal(wrapData)
		if err != nil {
			errors.SetError(w, err)
			return
		}
		w.Write(body)
	}
}

func (rw *Wrapper) Empty(ctrlFunc func(r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := ctrlFunc(r)
		if err != nil {
			errors.SetError(w, err)
			return
		}
		rw.response(w, map[string]interface{}{"message": "ok"})
	}
}

func (rw *Wrapper) Data(ctrlFunc func(r *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ctrlFunc(r)
		if err != nil {
			errors.SetError(w, err)
			return
		}
		rw.response(w, data)
	}
}
