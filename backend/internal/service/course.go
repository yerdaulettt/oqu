package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"

	"oqu/internal/models"
	"oqu/internal/repository"

	"github.com/redis/go-redis/v9"
)

type courseService struct {
	repo  repository.CourseRepository
	cache *redis.Client
}

func NewCourseService(r repository.CourseRepository, c *redis.Client) *courseService {
	return &courseService{repo: r, cache: c}
}

func (s *courseService) Get() ([]models.Course, error) {
	courses, err := s.repo.GetCourses()

	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Course service Get():", err)
		return nil, notFoundErr
	} else if err != nil {
		log.Println("Course service Get():", err)
		return nil, internalErr
	}

	return courses, nil
}

func (s *courseService) GetById(id int) (*models.Course, error) {
	ctx := context.Background()
	var course models.Course

	err := s.cache.HGetAll(ctx, "course-"+strconv.Itoa(id)).Scan(&course)
	if err != nil {
		c, err := s.repo.GetCourseById(id)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, notFoundErr
		} else if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		return c, nil
	}

	if course.Id == 0 {
		result, err := s.repo.GetCourseById(id)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, notFoundErr
		} else if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		key := "course-" + strconv.Itoa(id)
		err = s.cache.HSet(ctx, key, result).Err()
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		return result, nil
	}

	return &course, nil
}

func (s *courseService) GetCourseLessons(id int) ([]models.Lesson, error) {
	lessons, err := s.repo.GetCourseLessons(id)

	if err != nil {
		log.Println(err)
		return nil, internalErr
	} else if lessons == nil {
		return nil, notFoundErr
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
