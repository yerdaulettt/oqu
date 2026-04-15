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

func (s *lessonService) PostComment(lessonId int, userId int, c *models.Comment) bool {
	ok, err := s.repo.PostComment(lessonId, userId, c)
	if err != nil {
		log.Println(err)
		return false
	}
	return ok
}
