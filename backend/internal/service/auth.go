package service

import (
	"log"
	"oqu/internal/models"
	"oqu/internal/repository"
)

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) *authService {
	return &authService{repo: r}
}

func (s *authService) Register(u *models.UserRegister) int {
	id, err := s.repo.Register(u)
	if err != nil {
		log.Println(err)
		return -1
	}
	return id
}

func (s *authService) Login(u *models.UserLogin) string {
	token, err := s.repo.Login(u)
	if err != nil {
		return err.Error()
	}
	return token
}
