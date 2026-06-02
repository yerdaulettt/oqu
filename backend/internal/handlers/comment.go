package handlers

import (
	"log"
	"net/http"
	"oqu/internal/service"
	"strconv"
)

type commentHandler struct {
	srvc service.CommentService
}

func NewCommentHandler(s service.CommentService) *commentHandler {
	return &commentHandler{srvc: s}
}

func (h *commentHandler) Vote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "provide number")
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		w.Write([]byte(`{"error:" "unauthorized or bad request"}`))
		return
	}

	err = h.srvc.Vote(userId, commentId)
	if err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Write([]byte(`{"msg": "voted successfully"}`))
}

func (h *commentHandler) ModifyVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(`{"error": "provide number!"}`))
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		w.Write([]byte(`{"error": "bad request"}`))
		return
	}

	err = h.srvc.ModifyVote(userId, commentId)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	jsonResponse(w, http.StatusOK, "vote changed")
}
