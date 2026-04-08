package handlers

import (
	"encoding/json"
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

	token := h.srvc.Login(&userLog)
	w.Write([]byte(`{"access": "` + token + `"}`))
}
