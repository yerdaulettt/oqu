package handlers

import (
	"encoding/json"
	"net/http"

	"oqu/internal/service"
)

type courseHandler struct {
	service service.CourseService
}

func NewCourseHandler(s service.CourseService) *courseHandler {
	return &courseHandler{service: s}
}

func (c *courseHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	courses, err := c.service.Get()
	if err != nil {
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	err = json.NewEncoder(w).Encode(&courses)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
}
