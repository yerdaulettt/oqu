package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"oqu/internal/service"
)

type moderatorHandler struct {
	srvc service.ModeratorService
}

func NewModeratorHandler(s service.ModeratorService) *moderatorHandler {
	return &moderatorHandler{srvc: s}
}

// @Tags moderator
// @Produce json
// @Success 200 {array} models.ModeratorCommentView
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [moderator] can access! Your role is [admin or user]"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /moderator/comments [get]
func (h *moderatorHandler) ViewComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	comments, err := h.srvc.ViewComments()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&comments)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}

// @Tags moderator
// @Produce json
// @Param id path int true "Comment id"
// @Success 200 {object} models.DeletedComment
// @Failure 400 "Provide number"
// @Failure 401 "No token found, incorrect token or token expired"
// @Failure 403 "Only [moderator] can access! Your role is [admin or user]"
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Security Bearer
// @Router /moderator/comments/{id} [delete]
func (h *moderatorHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, numberErr.Error())
		return
	}

	deleted, err := h.srvc.DeleteComment(id)
	if errors.Is(err, service.NotFoundErr) {
		jsonResponse(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(&deleted)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, internalErr.Error())
	}
}
