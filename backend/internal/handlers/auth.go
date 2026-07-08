package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"oqu/internal/models"
	"oqu/internal/service"
)

type authHandler struct {
	srvc service.AuthService
}

func NewAuthHandler(s service.AuthService) *authHandler {
	return &authHandler{srvc: s}
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userReg models.UserRegister

	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	if len(userReg.Name) == 0 || len(userReg.Username) == 0 || len(userReg.Password) == 0 {
		jsonResponse(w, http.StatusBadRequest, "Provide correct info!")
		return
	}

	if len(userReg.Role) == 0 {
		userReg.Role = "user"
	}

	id := h.srvc.Register(&userReg)
	if id == -1 {
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	w.Write([]byte(`{"msg": "user with id ` + strconv.Itoa(id) + `"}`))
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userLog models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&userLog)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	token, err := h.srvc.Login(&userLog)
	if errors.Is(err, service.NotFoundErr) {
		jsonResponse(w, http.StatusNotFound, err.Error())
		return
	} else if errors.Is(err, service.IncorrectPassword) {
		jsonResponse(w, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write([]byte(`{"access": "` + token + `"}`))
}
