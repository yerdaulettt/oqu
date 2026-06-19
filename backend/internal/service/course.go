package service

import (
	"database/sql"
	"errors"
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

func (s *courseService) Get() ([]models.Course, error) {
	courses, err := s.repo.GetCourses()

	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Course service Get():", err)
		return nil, notFoundErr
	} else if err != nil {
		log.Println("Course service Get():", err)
		return nil, internalErr
	}

	return courses, nil
}

func (s *courseService) GetById(id int) (*models.Course, error) {
	course, err := s.repo.GetCourseById(id)

	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Course service GetById():", err)
		return nil, notFoundErr
	} else if err != nil {
		log.Println("Course service GetById():", err)
		return nil, internalErr
	}

	return course, nil
}

func (s *courseService) GetCourseLessons(id int) ([]models.Lesson, error) {
	lessons, err := s.repo.GetCourseLessons(id)

	if err != nil {
		log.Println(err)
		return nil, internalErr
	} else if lessons == nil {
		return nil, notFoundErr
	}

	return lessons, nil
}

func (s *courseService) EnrollInClass(classId int, userId int) error {
	err := s.repo.EnrollInClass(classId, userId)
	if err != nil {
		log.Println(err)
		return internalErr
	}
	return nil
}

func (s *courseService) ResetRating(courseId, userId int) error {
	err := s.repo.ResetRating(courseId, userId)
	if err != nil {
		log.Println(err)
		return internalErr
	}

	return nil
}
