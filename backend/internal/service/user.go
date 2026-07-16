package service

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"oqu/internal/models"
	"oqu/internal/repository"
)

type userService struct {
	repo  repository.UserRepository
	cache repository.CacheRepository
}

func NewUserService(r repository.UserRepository, c repository.CacheRepository) *userService {
	return &userService{repo: r, cache: c}
}

func (s *userService) GetProfileInfo(userId int) (*models.User, error) {
	ctx := context.Background()
	var user models.User
	key := "user-" + strconv.Itoa(userId)

	value, err := s.cache.Get(ctx, key)
	if err == nil && value != nil {
		if err := json.Unmarshal(value, &user); err == nil {
			return &user, nil
		}
	}

	profile, err := s.repo.GetProfileInfo(userId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	result, err := json.Marshal(profile)
	if err != nil {
		log.Println("marshal", err)
	}

	err = s.cache.Set(ctx, key, result, 5*time.Minute)
	if err != nil {
		log.Println("set", err)
	}

	return profile, nil
}

func (s *userService) UpdateProfile(u *models.UserUpdate, userId int) (*models.User, error) {
	params := []any{}
	columns := []string{}

	if u.Name != "" {
		columns = append(columns, "name")
		params = append(params, u.Name)
	}

	if u.Username != "" {
		exists, err := s.repo.UsernameExists(u.Username)
		if err != nil {
			log.Println(err)
			return nil, UsernameErr
		}

		if exists {
			return nil, UsernameErr
		}

		columns = append(columns, "username")
		params = append(params, u.Username)
	}

	params = append(params, userId)

	updated, err := s.repo.UpdateProfile(params, columns)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return updated, nil
}

func (s *userService) GetMyClasses(userId int) ([]models.Course, error) {
	courses, err := s.repo.GetMyClasses(userId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return courses, nil
}

func (s *userService) GetAllCoursesRating(userId int) ([]models.Rating, error) {
	rating, err := s.repo.GetAllCoursesRating(userId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return rating, nil
}
