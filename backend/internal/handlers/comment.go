package handlers

import (
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

	err = h.srvc.Vote(commentId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Write([]byte(`{"msg": "voted successfully"}`))
}
