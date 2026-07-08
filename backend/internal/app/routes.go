package app

import (
	"database/sql"
	"net/http"

	"oqu/internal/auth"
	"oqu/internal/handlers"
	"oqu/internal/middleware"
	"oqu/internal/repository"
	"oqu/internal/repository/postgresql"
	"oqu/internal/service"

	"github.com/go-chi/chi/v5"
)

func authRouter(db *sql.DB, jwtService *auth.JwtAuth) http.Handler {
	r := chi.NewRouter()

	authR := postgresql.NewAuthRepo(db)
	authS := service.NewAuthService(authR, jwtService)
	authH := handlers.NewAuthHandler(authS)

	r.Post("/register", authH.Register)
	r.Post("/login", authH.Login)

	return r
}

func courseRouter(db *sql.DB, cache repository.CacheRepository, jwtService *auth.JwtAuth) http.Handler {
	r := chi.NewRouter()

	courseR := postgresql.NewCourseRepo(db)
	courseS := service.NewCourseService(courseR, cache)
	courseH := handlers.NewCourseHandler(courseS)

	r.Use(middleware.JWTAuthMiddleware(jwtService))
	r.HandleFunc("GET /", courseH.Get)
	r.HandleFunc("GET /{id}", courseH.GetById)
	r.With(middleware.Role("user")).HandleFunc("POST /{id}/enroll", courseH.EnrollInClass)
	r.With(middleware.Role("user")).HandleFunc("POST /{id}/reset", courseH.ResetRating)

	return r
}

func lessonRouter(db *sql.DB, cache repository.CacheRepository, jwtService *auth.JwtAuth) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.JWTAuthMiddleware(jwtService))

	lessonR := postgresql.NewLessonRepo(db)
	lessonS := service.NewLessonService(lessonR, cache)
	lessonH := handlers.NewLessonHandler(lessonS)

	r.HandleFunc("GET /{id}", lessonH.GetLessonById)
	r.HandleFunc("GET /{id}/comments", lessonH.GetComments)
	r.HandleFunc("POST /{id}/comments", lessonH.PostComment)
	r.With(middleware.Role("user")).HandleFunc("POST /{id}/score", lessonH.Score)
	r.With(middleware.Role("user")).HandleFunc("POST /{id}/reset", lessonH.ResetScore)
	r.HandleFunc("GET /{id}/test", lessonH.GetTest)
	r.With(middleware.Role("user")).HandleFunc("POST /{id}/test", lessonH.SubmitTest)

	return r
}

func commentRouter(db *sql.DB, jwtService *auth.JwtAuth) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.JWTAuthMiddleware(jwtService))

	commentR := postgresql.NewCommentRepo(db)
	commentS := service.NewCommentService(commentR)
	commentH := handlers.NewCommentHandler(commentS)

	r.HandleFunc("POST /{id}/vote", commentH.Vote)
	r.HandleFunc("PATCH /{id}/vote", commentH.ModifyVote)

	return r
}

func adminRouter(db *sql.DB, jwtService *auth.JwtAuth) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.JWTAuthMiddleware(jwtService))
	r.Use(middleware.Role("admin"))

	adminR := postgresql.NewAdminRepo(db)
	adminS := service.NewAdminService(adminR)
	adminH := handlers.NewAdminHandler(adminS)

	r.HandleFunc("GET /users", adminH.GetUsers)
	r.HandleFunc("POST /courses", adminH.MakeCourse)
	r.HandleFunc("DELETE /courses/{id}", adminH.Delete)
	r.HandleFunc("POST /courses/{id}/lessons", adminH.AddLesson)
	r.HandleFunc("POST /lessons/{id}/test", adminH.AddTest)
	r.HandleFunc("GET /lessons/{id}/test", adminH.GetTest)

	return r
}

func moderatorRouter(db *sql.DB, jwtService *auth.JwtAuth) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.JWTAuthMiddleware(jwtService))
	r.Use(middleware.Role("moderator"))

	moderatorR := postgresql.NewModeratorRepo(db)
	moderatorS := service.NewModeratorService(moderatorR)
	moderatorH := handlers.NewModeratorHandler(moderatorS)

	r.HandleFunc("GET /comments", moderatorH.ViewComments)
	r.HandleFunc("DELETE /comments/{id}", moderatorH.DeleteComment)

	return r
}

func userRouter(db *sql.DB, cache repository.CacheRepository, jwtService *auth.JwtAuth) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.JWTAuthMiddleware(jwtService))

	userR := postgresql.NewUserRepo(db)
	userS := service.NewUserService(userR, cache)
	userH := handlers.NewUserHandler(userS)

	r.HandleFunc("GET /profile", userH.GetProfileInfo)
	r.With(middleware.Role("user")).HandleFunc("GET /enrollments", userH.GetMyClasses)
	r.With(middleware.Role("user")).HandleFunc("GET /rating", userH.GetAllCoursesRating)

	return r
}
