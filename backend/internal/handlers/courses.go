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

func (h *courseHandler) EnrollInClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "unauthorized or bad request")
		return
	}

	classId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "provide number")
		return
	}

	err = h.srvc.EnrollInClass(classId, userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, "Successfully enrolled!")
}
