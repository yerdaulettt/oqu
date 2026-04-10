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

func (s *adminService) MakeCourse(c *models.Course) int {
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
