package service

import (
	"oqu/internal/models"
	"oqu/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *userService {
	return &userService{repo: r}
}

func (s *userService) GetProfileInfo(userId int) (*models.User, error) {
	return s.repo.GetProfileInfo(userId)
}

func (s *userService) GetMyClasses(userId int) ([]models.Course, error) {
	return s.repo.GetMyClasses(userId)
}
