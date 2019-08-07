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

func jsonResponse(w http.ResponseWriter, r Response) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(r.Code)
	json.NewEncoder(w).Encode(r)
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
