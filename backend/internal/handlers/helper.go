package handlers

import (
	"errors"
	"net/http"
)

var (
	incorrectUserId = errors.New("Incorrect user id")
	internalErr     = errors.New("Internal server error")
	requestBodyErr  = errors.New("Incorrect request body")
	numberErr       = errors.New("Provide number")
)

func jsonResponse(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"msg": "` + msg + `"}`))
}
