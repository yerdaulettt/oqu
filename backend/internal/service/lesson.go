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

func (s *lessonService) GetLesson(id int) (*models.LessonDetail, error) {
	lesson, err := s.repo.GetLesson(id)

	if err != nil {
		log.Println(err)
		return nil, internalErr
	} else if lesson == nil {
		return nil, notFoundErr
	}

	return lesson, nil
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

func (s *lessonService) Score(lessonId, userId int) error {
	return s.repo.Score(lessonId, userId)
}

func (s *lessonService) ResetScore(lessonId, userId int) error {
	err := s.repo.ResetScore(lessonId, userId)
	if err != nil {
		log.Println(err)
		return internalErr
	}

	return nil
}
