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

// @Tags user
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 "Incorrect user id"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /api/users/profile [get]
func (h *userHandler) GetProfileInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, incorrectUserId.Error())
		return
	}

	profile, err := h.srvc.GetProfileInfo(userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&profile)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}

// @Tags user
// @Produce json
// @Success 200 {array} models.Course
// @Failure 400 "Incorrect user id"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [user] can access! Your role is [admin or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /api/users/enrollments [get]
func (h *userHandler) GetMyClasses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, incorrectUserId.Error())
		return
	}

	classes, err := h.srvc.GetMyClasses(userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&classes)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}

// @Tags user
// @Produce json
// @Success 200 {array} models.Rating
// @Failure 400 "Incorrect user id"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [user] can access! Your role is [admin or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /api/users/rating [get]
func (h *userHandler) GetAllCoursesRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, incorrectUserId.Error())
		return
	}

	ratings, err := h.srvc.GetAllCoursesRating(userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&ratings)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}
