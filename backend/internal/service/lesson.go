package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"sort"
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

func (s *lessonService) GetLesson(id, userId int) (*models.LessonDetail, error) {
	ctx := context.Background()
	var l models.LessonDetail
	key := "lesson-" + strconv.Itoa(id)

	value, err := s.cache.Get(ctx, key)
	if err == nil && value != nil {
		if err := json.Unmarshal(value, &l); err == nil {
			return &l, nil
		}
	}

	lesson, err := s.repo.GetLesson(id, userId)
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

func (s *lessonService) GetTest(lessonId, userId int) ([]models.StudentTestView, error) {
	ctx := context.Background()
	var lt []models.StudentTestView

	key := "lessontest-" + strconv.Itoa(lessonId)

	value, err := s.cache.Get(ctx, key)
	if err == nil && value != nil {
		if err := json.Unmarshal(value, &lt); err == nil {
			return lt, nil
		}
	}

	lessonTest, err := s.repo.GetTest(lessonId, userId)
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

func (s *lessonService) ResetTest(lessonId, userId int) error {
	err := s.repo.ResetTest(lessonId, userId)
	if err != nil {
		log.Println(err)
		return internalErr
	}

	return nil
}

func (s *lessonService) SubmitTest(lessonId, userId int, st []models.SubmitTest) (*models.ResultsOfTest, error) {
	isCompleted := s.repo.IsTestCompleted(lessonId, userId)
	if isCompleted {
		return nil, AlreadyCompleted
	}

	correctAnswers, err := s.repo.GetCorrectAnswers(lessonId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	totalQuestions := len(correctAnswers)

	if totalQuestions != len(st) {
		return nil, incorrectTestSubmit
	}

	sort.Slice(st, func(i, j int) bool {
		return st[i].QuestionId < st[j].QuestionId
	})

	points := 0
	for i := range totalQuestions {
		if correctAnswers[i].QuestionId == st[i].QuestionId && correctAnswers[i].CorrectChoice == st[i].SelectedChoice {
			points += 1
			correctAnswers[i].IsCorrect = true
		}

		correctAnswers[i].SelectedChoice = st[i].SelectedChoice
	}

	if completed := ((points * 100) / totalQuestions) > 90; completed {
		err = s.repo.SubmitTest(lessonId, userId, completed, st)
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}
	}

	return &models.ResultsOfTest{TotalQuestions: totalQuestions, Point: points, Results: correctAnswers}, nil
}
