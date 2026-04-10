package app

import (
	"database/sql"
	"net/http"

	"oqu/internal/handlers"
	"oqu/internal/middleware"
	"oqu/internal/repository/postgresql"
	"oqu/internal/service"

	"github.com/go-chi/chi/v5"
)

func authRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	authR := postgresql.NewAuthRepo(db)
	authS := service.NewAuthService(authR)
	authH := handlers.NewAuthHandler(authS)

	r.Post("/register", authH.Register)
	r.Post("/login", authH.Login)

	return r
}

func courseRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	courseR := postgresql.NewCourseRepo(db)
	courseS := service.NewCourseService(courseR)
	courseH := handlers.NewCourseHandler(courseS)

	r.HandleFunc("GET /", courseH.Get)
	r.HandleFunc("GET /{id}", courseH.GetById)
	r.HandleFunc("GET /{id}/lessons", courseH.GetCourseLessons)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)
		r.Use(middleware.Role("admin"))

		r.HandleFunc("POST /", courseH.MakeCourse)
		r.HandleFunc("DELETE /{id}", courseH.Delete)
	})

	return r
}

func lessonRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.JWTAuthMiddleware)
	r.Use(middleware.Role("admin", "user"))

	lessonR := postgresql.NewLessonRepo(db)
	lessonS := service.NewLessonService(lessonR)
	lessonH := handlers.NewLessonHandler(lessonS)

	r.HandleFunc("GET /{id}/comments", lessonH.GetComments)

	return r
}
