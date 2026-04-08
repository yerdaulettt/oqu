package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"oqu/internal/models"
	"oqu/internal/service"
)

type courseHandler struct {
	srvc service.CourseService
}

func NewCourseHandler(s service.CourseService) *courseHandler {
	return &courseHandler{srvc: s}
}

func (h *courseHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	courses := h.srvc.Get()
	if courses == nil {
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	err := json.NewEncoder(w).Encode(&courses)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
}

func (h *courseHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(`{"error": "provide number"}`))
		return
	}

	course := h.srvc.GetById(id)
	if course == nil {
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	err = json.NewEncoder(w).Encode(course)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
}

func (h *courseHandler) GetCourseLessons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(`{"error": "provide number"}`))
		return
	}

	lessons := h.srvc.GetCourseLessons(id)
	if lessons == nil {
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	err = json.NewEncoder(w).Encode(&lessons)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
}

func (h *courseHandler) MakeCourse(w http.ResponseWriter, r *http.Request) {
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

func (h *courseHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
