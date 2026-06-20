package service

import (
	"context"
	"log"
	"strconv"

	"oqu/internal/models"
	"oqu/internal/repository"

	"github.com/redis/go-redis/v9"
)

type userService struct {
	repo  repository.UserRepository
	cache *redis.Client
}

func NewUserService(r repository.UserRepository, c *redis.Client) *userService {
	return &userService{repo: r, cache: c}
}

func (s *userService) GetProfileInfo(userId int) (*models.User, error) {
	ctx := context.Background()
	var user models.User

	err := s.cache.HGetAll(ctx, "user-"+strconv.Itoa(userId)).Scan(&user)
	if err != nil {
		profile, err := s.repo.GetProfileInfo(userId)
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		return profile, nil
	}

	if user.Id == 0 {
		result, err := s.repo.GetProfileInfo(userId)

		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		key := "user-" + strconv.Itoa(userId)
		err = s.cache.HSet(ctx, key, result).Err()
		if err != nil {
			log.Println(err)
			return nil, internalErr
		}

		return result, nil
	}

	return &user, nil
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
