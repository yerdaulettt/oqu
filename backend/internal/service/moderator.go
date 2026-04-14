package service

import (
	"oqu/internal/models"
	"oqu/internal/repository"
)

type moderatorService struct {
	repo repository.ModeratorRepository
}

func NewModeratorService(r repository.ModeratorRepository) *moderatorService {
	return &moderatorService{repo: r}
}

func (s *moderatorService) ViewComments() ([]models.ModeratorCommentView, error) {
	return s.repo.ViewComments()
}

func (s *moderatorService) DeleteComment(id int) (*models.Comment, error) {
	return s.repo.DeleteComment(id)
}
