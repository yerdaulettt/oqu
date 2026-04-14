package handlers

import (
	"encoding/json"
	"net/http"
	"oqu/internal/service"
	"strconv"
)

type moderatorHandler struct {
	srvc service.ModeratorService
}

func NewModeratorHandler(s service.ModeratorService) *moderatorHandler {
	return &moderatorHandler{srvc: s}
}

func (h *moderatorHandler) ViewComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	comments, err := h.srvc.ViewComments()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&comments)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *moderatorHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "provide number")
		return
	}

	deleted, err := h.srvc.DeleteComment(id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	err = json.NewEncoder(w).Encode(&deleted)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
	}
}
