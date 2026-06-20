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

func (s *lessonService) GetTest(lessonId int) ([]models.StudentTestView, error) {
	tests, err := s.repo.GetTest(lessonId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return tests, nil
}

func (s *lessonService) SubmitTest(lessonId int, st []models.SubmitTest) (*models.ResultsOfTest, error) {
	correctAnswers, err := s.repo.GetCorrectAnswers(lessonId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	if len(correctAnswers) != len(st) {
		return nil, incorrectTestSubmit
	}

	totalQuestions := len(correctAnswers)
	points := 0
	for i := range totalQuestions {
		if correctAnswers[i].QuestionId == st[i].QuestionId && correctAnswers[i].CorrectOptionId == st[i].SelectedChoice {
			points += 1
		}
	}

	return &models.ResultsOfTest{TotalQuestions: totalQuestions, Point: points}, nil
}
