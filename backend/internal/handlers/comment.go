package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"oqu/internal/models"
	"oqu/internal/service"
)

type commentHandler struct {
	srvc service.CommentService
}

func NewCommentHandler(s service.CommentService) *commentHandler {
	return &commentHandler{srvc: s}
}

func (h *commentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, incorrectUserId.Error())
		return
	}

	var c models.NewComment
	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, requestBodyErr.Error())
		return
	}

	comment, err := h.srvc.UpdateComment(commentId, userId, c.Content)
	if err != nil {
		if errors.Is(err, service.AuthErr) {
			jsonResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, service.NotFoundErr) {
			jsonResponse(w, http.StatusNotFound, err.Error())
			return
		}

		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}
}

func (h *commentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, incorrectUserId.Error())
		return
	}

	comment, err := h.srvc.DeleteComment(commentId, userId)
	if err != nil {
		if errors.Is(err, service.AuthErr) {
			jsonResponse(w, http.StatusForbidden, err.Error())
			return
		}

		if errors.Is(err, service.NotFoundErr) {
			jsonResponse(w, http.StatusNotFound, err.Error())
			return
		}

		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}
}

func (h *commentHandler) Vote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		jsonResponse(w, http.StatusBadRequest, incorrectUserId.Error())
		return
	}

	vote, err := strconv.ParseBool(r.URL.Query().Get("value"))
	if err != nil {
		vote = true
	}

	err = h.srvc.Vote(commentId, userId, vote)
	if err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
		return
	}

	jsonResponse(w, http.StatusOK, "Vote changed")
}
