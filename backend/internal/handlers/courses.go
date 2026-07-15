package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"oqu/internal/service"
)

type courseHandler struct {
	srvc service.CourseService
}

func NewCourseHandler(s service.CourseService) *courseHandler {
	return &courseHandler{srvc: s}
}

// @Tags course
// @Produce json
// @Success 200 {array} models.Course
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /api/courses [get]
func (h *courseHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	courses, err := h.srvc.Get()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&courses)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}

// @Tags course
// @Produce json
// @Param id path int true "Course id"
// @Success 200 {array} models.CourseDetails
// @Failure 400 "Provide number"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /api/courses/{id} [get]
func (h *courseHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, incorrectUserId.Error())
		return
	}

	course, err := h.srvc.GetById(id, userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(course)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}

// @Tags course
// @Produce json
// @Param id path int true "Course id"
// @Success 200 "Successfully enrolled!"
// @Failure 400 "Incorrect user id or provide number error"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [user] can access! Your role is [admin or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /api/courses/{id}/enroll [post]
func (h *courseHandler) EnrollInClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, incorrectUserId.Error())
		return
	}

	classId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	err = h.srvc.EnrollInClass(classId, userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, "Successfully enrolled!")
}

// @Tags course
// @Produce json
// @Param id path int true "Course id"
// @Success 200 "Course rating reset"
// @Failure 400 "Incorrect user id or provide number error"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [user] can access! Your role is [admin or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /api/courses/{id}/reset [post]
func (h *courseHandler) ResetRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	courseId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, incorrectUserId.Error())
		return
	}

	err = h.srvc.ResetRating(courseId, userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, "Course rating reset")
}
