package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

var (
	errorRecordNotFound  = errors.New("no records found")
	errorCannotAddTeam   = errors.New("cannot add team to game")
	errorWrongGameStatus = errors.New("cannot update stats for this game")
	errorGameNotStarted  = errors.New("cannot change status of this game")
	errorSameTeam        = errors.New("teams cannot be the same")
	errorSharedMembers 			= errors.New("teams cannot have shared members")
	errorWrongNumPlayers = errors.New(fmt.Sprintf("team must contain between %d and %d members",
		gameMinNumMembers, gameMaxNumMembers))

	errorCodes = map[error]int{
	errorRecordNotFound:  http.StatusNotFound,
	errorCannotAddTeam:   http.StatusForbidden,
	errorWrongGameStatus: http.StatusForbidden,
	errorGameNotStarted:  http.StatusForbidden,
	errorSameTeam:        http.StatusForbidden,
	errorWrongNumPlayers: http.StatusConflict,
	errorSharedMembers:   http.StatusConflict,
	}
)

/*
	update response code and/or message if necessary
 */
func translateResponse(r *Response) {
	switch {
	case strings.Contains(r.Message, "UNIQUE constraint failed"):	/* gorm error */
		r.Message = "record already exists"
		r.Code = http.StatusNotAcceptable
	case r.Message == "EOF":												/* json decode error */
		r.Message = "unexpected end of data"
		r.Code = http.StatusBadRequest
	}
}

func responseJson(w http.ResponseWriter, e error, model interface{}, successCode int) {
	var res = Response{}
	if e == nil {
		res = Response{true, http.StatusOK, "ok", model}
		if successCode != 0 {
			res.Code = successCode
		}
	} else {
		res = Response{false, http.StatusInternalServerError, e.Error(), nil}
		if errorCodes[e] != 0 {
			res.Code = errorCodes[e]
		}
	}
	translateResponse(&res)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(res.Code)
	_ = json.NewEncoder(w).Encode(res)
}
