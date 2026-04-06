package service

import "oqu/internal/models"

type CourseService interface {
	Get() []models.Course
}
