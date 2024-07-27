package lib

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	SatatusCode int         `json:"status_code,omitempty"`
	Message     string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

func failResponseWriter(w http.ResponseWriter, err error, errStatusCode int) {
	w.Header().Set("Content-Type", "application/json")

	var resp Response
	w.WriteHeader(errStatusCode)
	resp.SatatusCode = errStatusCode
	resp.Message = err.Error()
	resp.Data = nil

	responseBytes, _ := json.Marshal(resp)
	w.Write(responseBytes)
}

func successResponseWriter(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	var resp Response
	w.WriteHeader(statusCode)
	resp.SatatusCode = statusCode
	resp.Message = "success"
	resp.Data = data

	responseBytes, _ := json.Marshal(resp)
	w.Write(responseBytes)
}

func WriteResponse(w http.ResponseWriter, err error, data any) {
	switch err.(type) {
	case *ErrForbidden, ErrForbidden:
		failResponseWriter(w, err, http.StatusForbidden)
	case *ErrUnauthorized, ErrUnauthorized:
		failResponseWriter(w, err, http.StatusUnauthorized)
	case *ErrNotFound, ErrNotFound:
		failResponseWriter(w, err, http.StatusNotFound)
	case *ErrBadRequest, ErrBadRequest:
		failResponseWriter(w, err, http.StatusBadRequest)
	case nil:
		successResponseWriter(w, data, http.StatusOK)
	default:
		failResponseWriter(w, err, http.StatusInternalServerError)
	}
}
