package app

import (
	"database/sql"
	"log"
	"net/http"

	"oqu/internal/handlers"
	"oqu/internal/middleware"
)

func Bastau(db *sql.DB) {
	r := http.NewServeMux()

	lesson := handlers.NewCourseHandler(db)

	r.HandleFunc("GET /courses", lesson.Get)

	log.Fatal(http.ListenAndServe(":8080", middleware.LogMiddleware(r)))
}
