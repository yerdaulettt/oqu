package service

import (
	"log"

	"oqu/internal/models"
	"oqu/internal/repository"
)

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(r repository.AdminRepository) *adminService {
	return &adminService{repo: r}
}

func (s *adminService) GetUsers() []models.User {
	users, err := s.repo.GetUsers()
	if err != nil {
		log.Println(err)
		return nil
	}
	return users
}

func (s *adminService) MakeCourse(c *models.NewCourse) int {
	id, err := s.repo.MakeCourse(c)
	if err != nil {
		log.Println(err)
		return 0
	}

	return id
}

func (s *adminService) Delete(id int) *models.Course {
	result, err := s.repo.DeleteCourse(id)
	if err != nil {
		log.Println("course service error:", err)
		return nil
	}

	return result
}

func (s *adminService) AddLesson(courseId int, l *models.NewLesson) (int, error) {
	id, err := s.repo.AddLesson(courseId, l)

	if err != nil {
		log.Println("Admin service AddLesson():", err)
		return id, internalErr
	}

	return id, nil
}

func (s *adminService) AddTest(lessonId int, t []*models.NewTest) error {
	err := s.repo.AddTest(lessonId, t)
	if err != nil {
		log.Println(err)
		return internalErr
	}

	return nil
}

func (s *adminService) GetTest(lessonId int) ([]models.AdminTestView, error) {
	tests, err := s.repo.GetTest(lessonId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return tests, nil
}

func (s *adminService) DeleteTest(lessonId int) error {
	err := s.repo.DeleteTest(lessonId)
	if err != nil {
		log.Println(err)
		return internalErr
	}

	return nil
}
