package usecase

import (
	"database/sql"
	"oqu/internal/repository"
)

type courseUsecase struct {
	repo repository.CourseRepository
}

func NewCourseUsecase(db *sql.DB) *courseUsecase {
	return &courseUsecase{
		repo: repository.NewRepositories(db),
	}
}

func (cu *courseUsecase) Get() {
	cu.repo.GetCourses()
}
