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

/*
	update error code based on message, update message based on code or both
 */
func translateResponse(r *Response) {
	switch {
	case strings.Contains(r.Message, "UNIQUE constraint failed"):
		r.Message = "record already exists"
		r.Code = http.StatusNotAcceptable
	case r.Message == "record not found":
		r.Code = http.StatusNotFound
	case strings.Contains(r.Message,"cannot add team to game"):
		r.Code = http.StatusForbidden
	case r.Message == "cannot update stats for this game":
		r.Code = http.StatusForbidden
	case r.Message == "cannot change status of this game":
		r.Code = http.StatusForbidden
	case r.Message == "teams cannot be the same":
		r.Code = http.StatusForbidden
	case r.Message == "teams must have same number of players":
		r.Code = http.StatusForbidden
	case r.Message == "teams cannot have shared members":
		r.Code = http.StatusForbidden

	}
}

func responseJson(w http.ResponseWriter, e error, model interface{}, successCode int) {
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