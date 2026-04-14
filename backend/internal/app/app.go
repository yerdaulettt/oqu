package app

import (
	"database/sql"
	"log"
	"net/http"

	"oqu/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func Bastau(db *sql.DB) {
	r := chi.NewRouter()

	r.Use(middleware.LogMiddleware)

	r.Mount("/auth", authRouter(db))
	r.Mount("/admin", adminRouter(db))
	r.Mount("/moderator", moderatorRouter(db))
	r.Mount("/api/courses", courseRouter(db))
	r.Mount("/api/lessons", lessonRouter(db))

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
