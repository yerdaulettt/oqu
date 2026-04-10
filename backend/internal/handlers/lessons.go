package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"oqu/internal/models"
	"oqu/internal/service"
)

type lessonHandler struct {
	srvc service.LessonService
}

func NewLessonHandler(s service.LessonService) *lessonHandler {
	return &lessonHandler{srvc: s}
}

func (h *lessonHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(`{"error": "provide number"}`))
		return
	}

	comments := h.srvc.GetComments(id)
	if comments == nil {
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	err = json.NewEncoder(w).Encode(&comments)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
}

func (h *lessonHandler) PostComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(`{"error: "enter number"}`))
		return
	}

	var c models.Comment
	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		w.Write([]byte(`{"error": "json error"}`))
		return
	}

	ok := h.srvc.PostComment(lessonId, &c)
	if !ok {
		w.Write([]byte(`{"error": "comment not posted"}`))
		return
	}

	w.Write([]byte(`{"msg": "posted"}`))
}
