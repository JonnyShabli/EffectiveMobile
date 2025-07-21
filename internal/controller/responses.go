package controller

import (
	"encoding/json"
	"net/http"
)

type response struct {
	msg    string
	data   interface{}
	errMsg string
}

func SuccessResponse(w http.ResponseWriter, msg string, data interface{}, errMsg string) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := response{
		msg:    msg,
		data:   data,
		errMsg: errMsg,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func BadRequestResponse(w http.ResponseWriter, msg string, data interface{}, errMsg string) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	resp := response{
		msg:    msg,
		data:   data,
		errMsg: errMsg,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func ErrorResponse(w http.ResponseWriter, msg string, data interface{}, errMsg string) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	resp := response{
		msg:    msg,
		data:   data,
		errMsg: errMsg,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
