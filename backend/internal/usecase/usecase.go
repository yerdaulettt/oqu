package usecase

import (
	"database/sql"
)

type CourseUsecase interface {
	Get()
}

type Usecases struct {
	CourseUsecase
}

func NewUsecases(db *sql.DB) *Usecases {
	return &Usecases{
		CourseUsecase: NewCourseUsecase(db),
	}
}
