package handlers

import (
	"encoding/json"
	"net/http"
	"oqu/internal/models"
	"oqu/internal/service"
	"strconv"
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

func (h *adminHandler) MakeCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c *models.Course
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	id := h.srvc.MakeCourse(c)
	if id == 0 {
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	w.Write([]byte(`{"msg": "course with id ` + strconv.Itoa(id) + `"}`))
}

func (h *adminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(`{"error": "provide number"}`))
		return
	}

	result := h.srvc.Delete(id)
	if result == nil {
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
}
