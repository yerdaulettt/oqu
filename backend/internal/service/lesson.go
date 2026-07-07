package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"oqu/internal/models"
	"oqu/internal/repository"
)

type lessonService struct {
	repo  repository.LessonRepository
	cache repository.CacheRepository
}

func NewLessonService(r repository.LessonRepository, c repository.CacheRepository) *lessonService {
	return &lessonService{repo: r, cache: c}
}

func (s *lessonService) GetLesson(id int) (*models.LessonDetail, error) {
	ctx := context.Background()
	var l models.LessonDetail
	key := "lesson-" + strconv.Itoa(id)

	value, err := s.cache.Get(ctx, key)
	if err == nil && value != nil {
		if err := json.Unmarshal(value, &l); err == nil {
			return &l, nil
		}
	}

	lesson, err := s.repo.GetLesson(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NotFoundErr
	} else if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	result, err := json.Marshal(lesson)
	if err != nil {
		log.Println("marshal", err)
	}

	err = s.cache.Set(ctx, key, result, 5*time.Minute)
	if err != nil {
		log.Println("set", err)
	}

	return lesson, nil
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
	var lt []models.StudentTestView

	key := "lessontest-" + strconv.Itoa(lessonId)

	value, err := s.cache.Get(ctx, key)
	if err == nil && value != nil {
		if err := json.Unmarshal(value, &lt); err == nil {
			return lt, nil
		}
	}

	lessonTest, err := s.repo.GetTest(lessonId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NotFoundErr
	} else if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	result, err := json.Marshal(lessonTest)
	if err != nil {
		log.Println("marshal", err)
	}

	err = s.cache.Set(ctx, key, result, 5*time.Minute)
	if err != nil {
		log.Println("set", err)
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
