package app

import (
	"database/sql"
	"log"
	"net/http"

	"oqu/internal/auth"
	"oqu/internal/middleware"
	"oqu/internal/repository"

	_ "oqu/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Bastau(db *sql.DB, cache repository.CacheRepository, jwtService *auth.JwtAuth) {
	r := chi.NewRouter()

	r.Use(middleware.LogMiddleware)

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	r.Mount("/auth", authRouter(db, jwtService))
	r.Mount("/admin", adminRouter(db, jwtService))
	r.Mount("/moderator", moderatorRouter(db, jwtService))
	r.Mount("/api/my", userRouter(db, cache, jwtService))
	r.Mount("/api/courses", courseRouter(db, cache, jwtService))
	r.Mount("/api/lessons", lessonRouter(db, cache, jwtService))
	r.Mount("/api/comments", commentRouter(db, jwtService))

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
