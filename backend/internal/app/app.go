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

	lesson := handlers.NewLessonHandler(db)

	r.HandleFunc("GET /lessons", lesson.Get)

	log.Fatal(http.ListenAndServe(":8080", middleware.LogMiddleware(r)))
}
