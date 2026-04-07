package app

import (
	"database/sql"
	"log"
	"net/http"

	"oqu/internal/handlers"
	"oqu/internal/middleware"
	"oqu/internal/repository/postgresql/course"
	"oqu/internal/service"
)

func Bastau(db *sql.DB) {
	r := http.NewServeMux()

	repo := course.NewCourseRepo(db)
	service := service.NewCourseService(repo)
	handler := handlers.NewCourseHandler(service)

	r.HandleFunc("GET /courses", handler.Get)
	r.HandleFunc("GET /courses/{id}", handler.GetById)
	r.HandleFunc("GET /courses/{id}/lessons", handler.GetCourseLessons)
	r.HandleFunc("POST /courses", handler.MakeCourse)
	r.HandleFunc("DELETE /courses/{id}", handler.Delete)

	log.Fatal(http.ListenAndServe(":8080", middleware.LogMiddleware(r)))
}
