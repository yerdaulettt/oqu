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

	switch status {
	case 400, 401, 402, 403, 404, 500:
		w.Write([]byte(`{"error": "` + msg + `"}`))
	default:
		w.Write([]byte(`{"message": "` + msg + `"}`))
	}
}
