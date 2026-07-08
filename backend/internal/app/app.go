package app

import (
	"database/sql"
	"log"
	"net/http"

	"oqu/internal/auth"
	"oqu/internal/middleware"
	"oqu/internal/repository"

	"github.com/go-chi/chi/v5"
)

func Bastau(db *sql.DB, cache repository.CacheRepository, jwtService *auth.JwtAuth) {
	r := chi.NewRouter()

	r.Use(middleware.LogMiddleware)

	r.Mount("/auth", authRouter(db, jwtService))
	r.Mount("/admin", adminRouter(db, jwtService))
	r.Mount("/moderator", moderatorRouter(db, jwtService))
	r.Mount("/api/users", userRouter(db, cache, jwtService))
	r.Mount("/api/courses", courseRouter(db, cache, jwtService))
	r.Mount("/api/lessons", lessonRouter(db, cache, jwtService))
	r.Mount("/api/comments", commentRouter(db, jwtService))

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
