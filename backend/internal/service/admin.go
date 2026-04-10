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
