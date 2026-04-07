package app

import (
	"database/sql"
	"log"
	"net/http"

	"oqu/internal/handlers"
	"oqu/internal/middleware"
	"oqu/internal/repository/postgresql/course"
	"oqu/internal/repository/postgresql/lesson"
	"oqu/internal/service"
)

func Bastau(db *sql.DB) {
	r := http.NewServeMux()

	courseR := course.NewCourseRepo(db)
	lessonR := lesson.NewLessonRepo(db)

	courseS := service.NewCourseService(courseR)
	lessonS := service.NewLessonService(lessonR)

	courseH := handlers.NewCourseHandler(courseS)
	lessonH := handlers.NewLessonHandler(lessonS)

	r.HandleFunc("GET /courses", courseH.Get)
	r.HandleFunc("GET /courses/{id}", courseH.GetById)
	r.HandleFunc("GET /courses/{id}/lessons", courseH.GetCourseLessons)
	r.HandleFunc("POST /courses", courseH.MakeCourse)
	r.HandleFunc("DELETE /courses/{id}", courseH.Delete)

	r.HandleFunc("GET /lessons/{id}/comments", lessonH.GetComments)

	log.Fatal(http.ListenAndServe(":8080", middleware.LogMiddleware(r)))
}
