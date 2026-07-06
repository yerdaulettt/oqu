package service

import (
	"context"
	"log"
	"strconv"

	"oqu/internal/models"
	"oqu/internal/repository"

	"github.com/redis/go-redis/v9"
)

type lessonService struct {
	repo  repository.LessonRepository
	cache *redis.Client
}

func NewLessonService(r repository.LessonRepository, c *redis.Client) *lessonService {
	return &lessonService{repo: r, cache: c}
}

func (s *lessonService) GetLesson(id int) (*models.LessonDetail, error) {
	ctx := context.Background()
	var lesson models.LessonDetail

	err := s.cache.HGetAll(ctx, "lesson-"+strconv.Itoa(id)).Scan(&lesson)
	if err != nil {
		lesson, err := s.repo.GetLesson(id)
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		return lesson, nil
	}

	if lesson.Id == 0 {
		result, err := s.repo.GetLesson(id)
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		key := "lesson-" + strconv.Itoa(id)
		err = s.cache.HSet(ctx, key, result).Err()
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		return result, nil
	}

	return &lesson, nil
}

func (s *lessonService) GetComments(lessonId, userId int) ([]models.Comment, error) {
	comments, err := s.repo.GetComments(lessonId, userId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return comments, nil
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
	err := s.repo.Score(lessonId, userId)
	if err != nil {
		log.Println(err)
		return internalErr
	}

	return nil
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
	ctx := context.Background()
	var lessonTest []models.StudentTestView

	err := s.cache.HGetAll(ctx, "lessonTest-"+strconv.Itoa(lessonId)).Scan(&lessonTest)
	if err != nil {
		test, err := s.repo.GetTest(lessonId)
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		return test, nil
	}

	if lessonTest == nil {
		result, err := s.repo.GetTest(lessonId)
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		key := "lessonTest-" + strconv.Itoa(lessonId)
		err = s.cache.HSet(ctx, key, result).Err()
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		return result, nil
	}

	return lessonTest, nil
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
