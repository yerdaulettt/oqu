package service

import (
	"oqu/internal/models"
	"oqu/internal/repository"
)

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(r repository.CourseRepository) *courseService {
	return &courseService{repo: r}
}

func (s *courseService) Get() ([]models.Course, error) {
	return s.repo.GetCourses()
}
