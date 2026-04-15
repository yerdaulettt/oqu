package service

import "oqu/internal/repository"

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(r repository.CommentRepository) *commentService {
	return &commentService{repo: r}
}

func (s *commentService) Vote(commentId int) error {
	return s.repo.Vote(commentId)
}
