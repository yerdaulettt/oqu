package service

import (
	"database/sql"
	"errors"
	"log"

	"oqu/internal/auth"
	"oqu/internal/models"
	"oqu/internal/repository"
)

type authService struct {
	repo       repository.AuthRepository
	jwtService *auth.JwtAuth
}

func NewAuthService(r repository.AuthRepository, j *auth.JwtAuth) *authService {
	return &authService{repo: r, jwtService: j}
}

func (s *authService) Register(u *models.UserRegister) int {
	hashed, err := hashPassword(u.Password)
	if err != nil {
		log.Println(err)
		return -1
	}

	u.Password = string(hashed)

	id, err := s.repo.Register(u)
	if err != nil {
		log.Println(err)
		return -1
	}

	return id
}

func (s *authService) Login(u *models.UserLogin) (string, error) {
	userFromDB, err := s.repo.GetUser(u.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", NotFoundErr
	} else if err != nil {
		log.Println(err)
		return "", internalErr
	}

	err = verifyPassword(userFromDB.PasswordHash, u.Password)
	if err != nil {
		return "", IncorrectPassword
	}

	tokenString, err := s.jwtService.GenerateTokens(userFromDB.Id, userFromDB.Role)
	if err != nil {
		log.Println(err)
		return "", internalErr
	}

	return tokenString.Access, nil
}
