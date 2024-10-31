package response

import (
	"encoding/json"
	"fmt"
	"github.com/mlplabs/app-utils/pkg/http/errors"
	"net/http"
	"reflect"
)

type PlainData struct {
	Data interface{} `json:"data"`
}

type List struct {
	Data  interface{} `json:"data"`
	Count int         `json:"count"`
}

type DataRange struct {
	Count  int32 `json:"count"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type Pagination struct {
	Data interface{} `json:"data"`
	DataRange
}

type Wrapper struct{}

func NewWrapper() *Wrapper {
	return &Wrapper{}
}

func (rw *Wrapper) response(w http.ResponseWriter, data interface{}) {
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			errors.SetError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
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

func (rw *Wrapper) DataPlain(ctrlFunc func(r *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ctrlFunc(r)
		if err != nil {
			errors.SetError(w, err)
			return
		}
		rw.response(w, PlainData{
			Data: data,
		})
	}
}

func (rw *Wrapper) DataList(ctrlFunc func(r *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ctrlFunc(r)
		if err != nil {
			errors.SetError(w, err)
			return
		}
		var listCount int
		switch reflect.TypeOf(data).Kind() {
		case reflect.Slice:
			listCount = reflect.ValueOf(data).Len()
		default:
			panic("return data does not common")
		}
		rw.response(w, List{
			Data:  data,
			Count: listCount,
		})
	}
}

func (rw *Wrapper) DataPages(ctrlFunc func(r *http.Request) (interface{}, *DataRange, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, params, err := ctrlFunc(r)
		if err != nil {
			errors.SetError(w, err)
			return
		}
		rw.response(w, Pagination{
			Data:      data,
			DataRange: *params,
		})
	}
}

func (rw *Wrapper) Raw(ctrlFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ctrlFunc(w, r)
		if err != nil {
			errors.SetError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, data)
	}
}
