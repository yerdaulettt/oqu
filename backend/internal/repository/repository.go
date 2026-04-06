package repository

import "oqu/internal/models"

type CourseRepository interface {
	GetCourses() ([]models.Course, error)
	GetCourseById(id int) (*models.Course, error)
	DeleteCourse(id int) (*models.Course, error)
}
