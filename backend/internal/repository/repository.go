package repository

import "oqu/internal/models"

type CourseRepository interface {
	GetCourses() ([]models.Course, error)
}
