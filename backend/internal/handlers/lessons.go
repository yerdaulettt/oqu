package handlers

import (
	"database/sql"
	"net/http"
)

type lessonHandler struct {
	db *sql.DB
}

func NewLessonHandler(db *sql.DB) *lessonHandler {
	return &lessonHandler{
		db: db,
	}
}

func (l *lessonHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"msg":"all lessons"}`))
}
