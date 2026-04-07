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

func (s *courseService) GetById(id int) *models.Course {
	course, err := s.repo.GetCourseById(id)
	if err != nil {
		log.Println("course service error:", err)
		return nil
	}

	return course
}

func (s *courseService) GetCourseLessons(id int) []models.Lesson {
	lessons, err := s.repo.GetCourseLessons(id)
	if err != nil {
		log.Println(err)
		return nil
	}

	return lessons
}

func (s *courseService) MakeCourse(c *models.Course) int {
	id, err := s.repo.MakeCourse(c)
	if err != nil {
		log.Println(err)
		return 0
	}

	return id
}

func (s *courseService) Delete(id int) *models.Course {
	result, err := s.repo.DeleteCourse(id)
	if err != nil {
		log.Println("course service error:", err)
		return nil
	}

	return result
}
