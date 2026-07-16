package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"oqu/internal/models"
	"oqu/internal/service"
)

type adminHandler struct {
	srvc service.AdminService
}

func NewAdminHandler(s service.AdminService) *adminHandler {
	return &adminHandler{srvc: s}
}

// @Tags admin
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [admin] can access! Your role is [user or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /admin/users [get]
func (h *adminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := h.srvc.GetUsers()
	if users == nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	err := json.NewEncoder(w).Encode(&users)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}

func (h *adminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	role := r.URL.Query().Get("role")
	if role == "" {
		jsonResponse(w, http.StatusBadRequest, "Role not found")
		return
	}

	user, err := h.srvc.UpdateUserRole(userId, role)
	if err != nil {
		if errors.Is(err, service.NotFoundErr) {
			jsonResponse(w, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, service.IncorrectRole) {
			jsonResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}
}

func (h *adminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	user, err := h.srvc.DeleteUser(userId)
	if err != nil {
		if errors.Is(err, service.NotFoundErr) {
			jsonResponse(w, http.StatusNotFound, err.Error())
			return
		}

		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}
}

// @Tags admin
// @Accept json
// @Produce json
// @Param course body models.NewCourse true "New course"
// @Success 200 "Course with id 1"
// @Failure 400 "Incorrect request body"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [admin] can access! Your role is [user or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /admin/courses [post]
func (h *adminHandler) MakeCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c *models.NewCourse
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	id := h.srvc.MakeCourse(c)
	if id == 0 {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	jsonResponse(w, http.StatusOK, "course with id"+strconv.Itoa(id))
}

func (h *adminHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	courseId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	var c models.NewCourse
	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	updated, err := h.srvc.UpdateCourse(&c, courseId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	err = json.NewEncoder(w).Encode(updated)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}
}

// @Tags admin
// @Produce json
// @Param id path int true "Course id"
// @Success 200 {object} models.Course
// @Failure 400 "Provide number"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [admin] can access! Your role is [user or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /admin/courses/{id} [delete]
func (h *adminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	result := h.srvc.Delete(id)
	if result == nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}

// @Tags admin
// @Accept json
// @Produce json
// @Param lesson body models.NewLesson true "New lesson"
// @Param id path int true "Course id"
// @Success 200 "Lesson with id 1"
// @Failure 400 "Incorrect request body"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [admin] can access! Your role is [user or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /admin/courses/{id}/lessons [post]
func (h *adminHandler) AddLesson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	courseId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, numberErr.Error())
		return
	}

	var l *models.NewLesson
	err = json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	id, err := h.srvc.AddLesson(courseId, l)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, "lesson with id "+strconv.Itoa(id))
}

func (h *adminHandler) UpdateLesson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	var l models.NewLesson
	err = json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	lesson, err := h.srvc.UpdateLesson(lessonId, &l)
	if err != nil {
		if errors.Is(err, service.NotFoundErr) {
			jsonResponse(w, http.StatusNotFound, err.Error())
			return
		}

		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	err = json.NewEncoder(w).Encode(lesson)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}
}

func (h *adminHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	lesson, err := h.srvc.DeleteLesson(lessonId)
	if err != nil {
		if errors.Is(err, service.NotFoundErr) {
			jsonResponse(w, http.StatusNotFound, err.Error())
			return
		}

		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	err = json.NewEncoder(w).Encode(lesson)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}
}

// @Tags admin
// @Accept json
// @Produce json
// @Param test body models.NewTest true "New test"
// @Param id path int true "Lesson id"
// @Success 201 "New test"
// @Failure 400 "Json error"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [admin] can access! Your role is [user or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /admin/lessons/{id}/test [post]
func (h *adminHandler) AddTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	var t []*models.NewTest
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	err = h.srvc.AddTest(lessonId, t)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, "New test")
}

// @Tags admin
// @Produce json
// @Param id path int true "Lesson id"
// @Success 200 {object} models.AdminTestView
// @Failure 400 "Provide number"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [admin] can access! Your role is [user or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /admin/lessons/{id}/test [get]
func (h *adminHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	tests, err := h.srvc.GetTest(lessonId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&tests)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}
}

// @Tags admin
// @Produce json
// @Param id path int true "Lesson id"
// @Success 200 "Deleted"
// @Failure 400 "Provide number"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [admin] can access! Your role is [user or moderator]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /admin/lessons/{id}/test [delete]
func (h *adminHandler) DeleteTest(w http.ResponseWriter, r *http.Request) {
	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	err = h.srvc.DeleteTest(lessonId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, "Deleted")
}
