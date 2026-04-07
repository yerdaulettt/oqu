package service

import "oqu/internal/models"

type CourseService interface {
	Get() []models.Course
	GetById(id int) *models.Course
	GetCourseLessons(id int) []models.Lesson
	Delete(id int) *models.Course
}
