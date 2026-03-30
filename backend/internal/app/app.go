package app

import (
	"database/sql"
	"log"
	"net/http"

	"oqu/internal/handlers"
)

func Bastau(db *sql.DB) {
	r := http.NewServeMux()

	lesson := handlers.NewLessonHandler(db)

	r.HandleFunc("GET /lessons", lesson.Get)

	log.Fatal(http.ListenAndServe(":8080", r))
}
