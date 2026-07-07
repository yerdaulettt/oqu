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

type courseService struct {
	repo  repository.CourseRepository
	cache repository.CacheRepository
}

func NewCourseService(r repository.CourseRepository, c repository.CacheRepository) *courseService {
	return &courseService{repo: r, cache: c}
}

func (s *courseService) Get() ([]models.Course, error) {
	courses, err := s.repo.GetCourses()

	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Course service Get():", err)
		return nil, NotFoundErr
	} else if err != nil {
		log.Println("Course service Get():", err)
		return nil, internalErr
	}

	return courses, nil
}

func (s *courseService) GetById(id int) (*models.Course, error) {
	ctx := context.Background()
	var c models.Course
	key := "course-" + strconv.Itoa(id)

	value, err := s.cache.Get(ctx, key)
	if err == nil && value != nil {
		if err := json.Unmarshal(value, &c); err == nil {
			return &c, nil
		}
	}

	course, err := s.repo.GetCourseById(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NotFoundErr
	} else if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	result, err := json.Marshal(course)
	if err != nil {
		log.Println("marshal", err)
	}

	err = s.cache.Set(ctx, key, result, 5*time.Minute)
	if err != nil {
		log.Println("set", err)
	}

	return course, nil
}

func (s *courseService) GetCourseLessons(id int) ([]models.Lesson, error) {
	lessons, err := s.repo.GetCourseLessons(id)

	if err != nil {
		log.Println(err)
		return nil, internalErr
	} else if lessons == nil {
		return nil, NotFoundErr
	}

	return lessons, nil
}

func (s *courseService) EnrollInClass(classId int, userId int) error {
	err := s.repo.EnrollInClass(classId, userId)
	if err != nil {
		log.Println(err)
		return internalErr
	}
	return nil
}

func (s *courseService) ResetRating(courseId, userId int) error {
	err := s.repo.ResetRating(courseId, userId)
	if err != nil {
		log.Println(err)
		return internalErr
	}

	return nil
}
