package handlers

import (
	"encoding/json"
	"net/http"
	"oqu/internal/service"
	"strconv"
)

type lessonHandler struct {
	srv service.LessonService
}

func NewLessonHandler(s service.LessonService) *lessonHandler {
	return &lessonHandler{srv: s}
}

func (h *lessonHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(`{"error": "provide number"}`))
		return
	}

	comments := h.srv.GetComments(id)
	if comments == nil {
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	err = json.NewEncoder(w).Encode(&comments)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
}
