package handlers

import (
	"encoding/json"
	"net/http"
	"oqu/internal/service"
)

type userHandler struct {
	srvc service.UserService
}

func NewUserHandler(s service.UserService) *userHandler {
	return &userHandler{srvc: s}
}

func (h *userHandler) GetProfileInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "unauthorized or bad request")
		return
	}

	profile, err := h.srvc.GetProfileInfo(userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	err = json.NewEncoder(w).Encode(&profile)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *userHandler) GetMyClasses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	classes, err := h.srvc.GetMyClasses(userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&classes)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *userHandler) GetAllCoursesRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "jwt claim error")
		return
	}

	ratings, err := h.srvc.GetAllCoursesRating(userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	err = json.NewEncoder(w).Encode(&ratings)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
	}
}
