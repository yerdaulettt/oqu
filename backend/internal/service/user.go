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
		log.Println(err)
		log.Println("CACHE ERROR GOING TO PG")
		return s.repo.GetProfileInfo(userId)
	}

	if user.Id == 0 {
		result, err := s.repo.GetProfileInfo(userId)

		if err != nil {
			return nil, err
		}

		key := "user-" + strconv.Itoa(userId)
		err = s.cache.HSet(ctx, key, result).Err()
		if err != nil {
			log.Println("REDIS ERROR")
			return nil, err
		}

		log.Println("NO CACHE, PG RESULT, CACHE SET")
		return result, nil
	}

	log.Println("FROM CACHE")
	return &user, nil
}

func (s *userService) GetMyClasses(userId int) ([]models.Course, error) {
	return s.repo.GetMyClasses(userId)
}

func (s *userService) GetAllCoursesRating(userId int) ([]models.Rating, error) {
	return s.repo.GetAllCoursesRating(userId)
}
