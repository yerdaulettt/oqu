package service

import (
	"log"
	"oqu/internal/models"
	"oqu/internal/repository"
)

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(r repository.CourseRepository) *courseService {
	return &courseService{repo: r}
}

func (s *courseService) Get() []models.Course {
	courses, err := s.repo.GetCourses()
	if err != nil {
		log.Println("course service error:", err)
		return nil
	}

	return courses
}
