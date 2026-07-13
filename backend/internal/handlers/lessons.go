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

func (h *lessonHandler) GetLessonById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Provide number")
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "incorrect user id")
		return
	}

	lesson, err := h.srvc.GetLesson(id, userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&lesson)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *lessonHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Provide number")
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "Incorrect user id")
		return
	}

	comments, err := h.srvc.GetComments(lessonId, userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&comments)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
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

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "error jwt claim")
		return
	}

	ok = h.srvc.PostComment(lessonId, userId, &c)
	if !ok {
		w.Write([]byte(`{"error": "comment not posted"}`))
		return
	}

	w.Write([]byte(`{"msg": "posted"}`))
}

func (h *lessonHandler) Score(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "provide number")
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "jwt claim error")
		return
	}

	response := h.srvc.Score(lessonId, userId)
	if response != nil {
		jsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Write([]byte(`{"msg": "scored ok"}`))
}

func (h *lessonHandler) ResetScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Provide number")
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "Incorrect token")
		return
	}

	err = h.srvc.ResetScore(lessonId, userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write([]byte(`{"message": "Score reset ok"}`))
}

func (h *lessonHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Provide number")
		return
	}

	test, err := h.srvc.GetTest(lessonId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&test)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *lessonHandler) SubmitTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Provide number")
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, "incorrect user")
		return
	}

	var st []models.SubmitTest
	err = json.NewDecoder(r.Body).Decode(&st)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.srvc.SubmitTest(lessonId, userId, st)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&result)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
