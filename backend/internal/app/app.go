package app

import (
	"database/sql"
	"log"
	"net/http"

	"oqu/internal/middleware"
	"oqu/internal/repository"

	"github.com/go-chi/chi/v5"
)

func Bastau(db *sql.DB, cache repository.CacheRepository) {
	r := chi.NewRouter()

	r.Use(middleware.LogMiddleware)

	r.Mount("/auth", authRouter(db))
	r.Mount("/admin", adminRouter(db))
	r.Mount("/moderator", moderatorRouter(db))
	r.Mount("/api/users", userRouter(db, cache))
	r.Mount("/api/courses", courseRouter(db, cache))
	r.Mount("/api/lessons", lessonRouter(db, cache))
	r.Mount("/api/comments", commentRouter(db))

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
