package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ResponseERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		ResponseJSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	ResponseJSON(w, http.StatusBadRequest, nil)
}

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Vary", "Accept-Encoding")
		next(w, r)
	}
}
