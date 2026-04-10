package handlers

import (
	"encoding/json"
	"net/http"
	"oqu/internal/service"
)

type adminHandler struct {
	srvc service.AdminService
}

func NewAdminHandler(s service.AdminService) *adminHandler {
	return &adminHandler{srvc: s}
}

func (h *adminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := h.srvc.GetUsers()
	if users == nil {
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	err := json.NewEncoder(w).Encode(&users)
	if err != nil {
		w.Write([]byte(`{"error": "json error"}`))
	}
}
