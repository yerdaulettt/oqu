package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"oqu/internal/auth"
	"oqu/internal/models"
	"oqu/internal/service"
)

type authHandler struct {
	srvc service.AuthService
}

func NewAuthHandler(s service.AuthService) *authHandler {
	return &authHandler{srvc: s}
}

// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRegister true "User registration"
// @Success 201 "User with id 1"
// @Failure 400 "Incorrect request body"
// @Failure 500 "Internal server error"
// @Router /auth/register [post]
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userReg models.UserRegister

	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	if len(userReg.Name) == 0 || len(userReg.Username) == 0 || len(userReg.Password) == 0 {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	if len(userReg.Role) == 0 {
		userReg.Role = "user"
	}

	id := h.srvc.Register(&userReg)
	if id == -1 {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, "User with id "+strconv.Itoa(id))
}

// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserLogin true "Login"
// @Success 200 {object} models.Tokens
// @Failure 400 "Incorrect request body"
// @Failure 403 "Incorrect password"
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /auth/login [post]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userLog models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&userLog)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	tokens, err := h.srvc.Login(&userLog)
	if errors.Is(err, service.NotFoundErr) {
		jsonResponse(w, http.StatusNotFound, err.Error())
		return
	} else if errors.Is(err, service.IncorrectPassword) {
		jsonResponse(w, http.StatusForbidden, err.Error())
		return
	} else if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}

// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body models.RefreshRequest true "Token refresh"
// @Success 200 "Access token"
// @Failure 400 "Incorrect request body or no refresh token"
// @Failure 403 "Token expired"
// @Failure 500 "Internal server error"
// @Router /auth/refresh [post]
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var refresh models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&refresh); err != nil {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	if refresh.Refresh == "" {
		jsonResponse(w, http.StatusBadRequest, "No refresh token")
		return
	}

	access, err := h.srvc.Refresh(refresh.Refresh)
	if errors.Is(err, auth.IncorrectToken) {
		jsonResponse(w, http.StatusBadRequest, err.Error())
		return
	} else if errors.Is(err, auth.ExpiredErr) {
		jsonResponse(w, http.StatusForbidden, err.Error())
		return
	} else if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	w.Write([]byte(`{"access": "` + access + `"}`))
}
