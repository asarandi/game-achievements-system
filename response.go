package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}


func translateError(code int, msg string) (int, string) {
	switch {
	case strings.Contains(msg, "UNIQUE constraint failed"):
		msg = "record already exists"
		code = http.StatusNotAcceptable
	case msg == "record not found":
		code = http.StatusNotFound
	}
	return code, msg
}

func jsonResponse(w http.ResponseWriter, r *Response) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(r.Code)
	_ = json.NewEncoder(w).Encode(*r)
}

func translateResponse(r *Response) {
	switch {
	case strings.Contains(r.Message, "UNIQUE constraint failed"):
		r.Message = "record already exists"
		r.Code = http.StatusNotAcceptable
	case r.Message == "record not found":
		r.Code = http.StatusNotFound
	case r.Message == "cannot update stats for this game":
		r.Code = http.StatusForbidden
	}
}

func responseDb(w http.ResponseWriter, e error, model interface{}, successCode int) {
	var res Response
	if successCode == 0 {
		successCode = http.StatusOK
	}
	if e != nil {
		res = Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
			Result:  nil,
		}
	} else {
		res = Response{
			Success: true,
			Code:    successCode,
			Message: "ok",
			Result:  model,
		}
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(res.Code)
	translateResponse(&res)
	_ = json.NewEncoder(w).Encode(res)
}