package course

import (
	"database/sql"
	"fmt"
)

type courseRepo struct {
	db *sql.DB
}

func NewCourseRepo(db *sql.DB) *courseRepo {
	return &courseRepo{db: db}
}

func (c *courseRepo) GetCourses() {
	fmt.Println("All courses")
}
