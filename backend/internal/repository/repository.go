package repository

import (
	"database/sql"

	"oqu/internal/repository/postgresql/course"
)

type CourseRepository interface {
	GetCourses()
}

type repositories struct {
	CourseRepository
}

func NewRepositories(db *sql.DB) *repositories {
	return &repositories{
		CourseRepository: course.NewCourseRepo(db),
	}
}
