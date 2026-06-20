package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"oqu/internal/models"
	"oqu/internal/service"
	"strconv"
)

type adminHandler struct {
	srvc service.AdminService
}

func NewAdminHandler(s service.AdminService) *adminHandler {
	return &adminHandler{srvc: s}
}

func (h *adminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := h.srvc.GetUsers()
	if users == nil {
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	err := json.NewEncoder(w).Encode(&users)
	if err != nil {
		w.Write([]byte(`{"error": "json error"}`))
	}
}

func (h *adminHandler) MakeCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c *models.Course
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	id := h.srvc.MakeCourse(c)
	if id == 0 {
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	w.Write([]byte(`{"msg": "course with id ` + strconv.Itoa(id) + `"}`))
}

func (h *adminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(`{"error": "provide number"}`))
		return
	}

	result := h.srvc.Delete(id)
	if result == nil {
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
}

func (h *adminHandler) AddLesson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	courseId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(`{"error": "provide number"}`))
		return
	}

	var l *models.Lesson
	err = json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		log.Println(err)
		w.Write([]byte(`{"error": "json error"}`))
		return
	}

	id, err := h.srvc.AddLesson(courseId, l)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"msg": "lesson with id ` + strconv.Itoa(id) + `"}`))
}

func (h *adminHandler) AddTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Provide number")
		return
	}

	var t []*models.NewTest
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.srvc.AddTest(lessonId, t)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write([]byte(`{"message": "New test"}`))
}

func (h *adminHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Provide number")
		return
	}

	tests, err := h.srvc.GetTest(lessonId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&tests)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
