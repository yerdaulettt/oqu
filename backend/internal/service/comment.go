package service

import (
	"log"
	"oqu/internal/repository"
)

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(r repository.CommentRepository) *commentService {
	return &commentService{repo: r}
}

func (s *commentService) Vote(userId, commentId int) error {
	return s.repo.Vote(userId, commentId)
}

func (s *commentService) ModifyVote(userId, commentId int) error {
	err := s.repo.ModifyVote(userId, commentId)
	if err != nil {
		log.Println("Comment service ModifyVote():", err)
		return internalErr
	}

	return nil
}
