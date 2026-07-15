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

func (s *authService) Login(u *models.UserLogin) (*models.Tokens, error) {
	userFromDB, err := s.repo.GetUser(u.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NotFoundErr
	} else if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	err = verifyPassword(userFromDB.PasswordHash, u.Password)
	if err != nil {
		return nil, IncorrectPassword
	}

	tokens, err := s.jwtService.GenerateTokens(userFromDB.Id, userFromDB.Role)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return tokens, nil
}

func (s *authService) Refresh(refresh string) (string, error) {
	access, err := s.jwtService.RefreshAccessToken(refresh)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return access, nil
}
