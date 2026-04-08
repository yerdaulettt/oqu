package app

import (
	"database/sql"
	"log"
	"net/http"

	"oqu/internal/handlers"
	"oqu/internal/middleware"
	"oqu/internal/repository/postgresql"
	"oqu/internal/service"
)

func Bastau(db *sql.DB) {
	r := http.NewServeMux()

	authR := postgresql.NewAuthRepo(db)
	courseR := postgresql.NewCourseRepo(db)
	lessonR := postgresql.NewLessonRepo(db)

	authS := service.NewAuthService(authR)
	courseS := service.NewCourseService(courseR)
	lessonS := service.NewLessonService(lessonR)

	authH := handlers.NewAuthHandler(authS)
	courseH := handlers.NewCourseHandler(courseS)
	lessonH := handlers.NewLessonHandler(lessonS)

	r.HandleFunc("POST /register", authH.Register)
	r.HandleFunc("POST /login", authH.Login)

	r.Handle("GET /courses", middleware.JWTAuthMiddleware(http.HandlerFunc(courseH.Get))) // auth test
	// r.HandleFunc("GET /courses", courseH.Get)
	r.HandleFunc("GET /courses/{id}", courseH.GetById)
	r.HandleFunc("GET /courses/{id}/lessons", courseH.GetCourseLessons)
	r.HandleFunc("POST /courses", courseH.MakeCourse)
	r.HandleFunc("DELETE /courses/{id}", courseH.Delete)

	r.HandleFunc("GET /lessons/{id}/comments", lessonH.GetComments)

	log.Fatal(http.ListenAndServe(":8080", middleware.LogMiddleware(r)))
}
