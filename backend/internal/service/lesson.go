package service

import (
	"log"
	"oqu/internal/models"
	"oqu/internal/repository"
)

type lessonService struct {
	repo repository.LessonRepository
}

func NewLessonService(r repository.LessonRepository) *lessonService {
	return &lessonService{repo: r}
}

func (s *lessonService) GetComments(id int) []models.Comment {
	comments, err := s.repo.GetComments(id)
	if err != nil {
		log.Println(err)
		return nil
	}

	return comments
}

func (s *lessonService) PostComment(lessonId int, c *models.Comment) bool {
	ok, err := s.repo.PostComment(lessonId, c)
	if err != nil {
		log.Println(err)
		return false
	}
	return ok
}
