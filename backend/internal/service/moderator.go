package service

import (
	"database/sql"
	"errors"
	"log"

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
	comments, err := s.repo.ViewComments()
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return comments, nil
}

func (s *moderatorService) DeleteComment(id int) (*models.Comment, error) {
	comment, err := s.repo.DeleteComment(id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, NotFoundErr
	} else if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return comment, nil
}
