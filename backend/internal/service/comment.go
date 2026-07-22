package service

import (
	"database/sql"
	"errors"
	"log"

	"oqu/internal/models"
	"oqu/internal/repository"
)

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(r repository.CommentRepository) *commentService {
	return &commentService{repo: r}
}

func (s *commentService) UpdateComment(commentId, userId int, content string) (*models.UpdatedComment, error) {
	userIdComment, err := s.repo.GetUserId(commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErr
		}

		log.Println(err)
		return nil, internalErr
	}

	if userId != userIdComment {
		return nil, AuthErr
	}

	comment, err := s.repo.UpdateComment(commentId, userId, content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErr
		}

		log.Println(err)
		return nil, internalErr
	}

	return comment, nil
}

func (s *commentService) DeleteComment(commentId, userId int) (*models.DeletedComment, error) {
	userIdComment, err := s.repo.GetUserId(commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErr
		}

		log.Println(err)
		return nil, internalErr
	}

	if userId != userIdComment {
		return nil, AuthErr
	}

	comment, err := s.repo.DeleteComment(commentId, userId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return comment, nil
}

func (s *commentService) Vote(commentId, userId int, vote bool) error {
	err := s.repo.Vote(commentId, userId, vote)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NotFoundErr
		}

		log.Println(err)
		return internalErr
	}

	return nil
}
