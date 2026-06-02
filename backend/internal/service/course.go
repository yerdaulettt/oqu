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

func (s *courseService) GetCourseLessons(id int) []models.Lesson {
	lessons, err := s.repo.GetCourseLessons(id)
	if err != nil {
		log.Println(err)
		return nil
	}

	return lessons
}

func (s *courseService) EnrollInClass(classId int, userId int) error {
	return s.repo.EnrollInClass(classId, userId)
}
