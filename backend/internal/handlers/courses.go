package handlers

import (
	"database/sql"
	"net/http"

	"oqu/internal/usecase"
)

type courseHandler struct {
	usecase usecase.CourseUsecase
}

func NewCourseHandler(db *sql.DB) *courseHandler {
	return &courseHandler{
		usecase: usecase.NewCourseUsecase(db),
	}
}

func (c *courseHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	c.usecase.Get()
	w.Write([]byte(`{"msg":"all courses"}`))
}
