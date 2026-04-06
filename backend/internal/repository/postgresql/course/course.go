package course

import (
	"database/sql"
	"oqu/internal/models"
)

type courseRepo struct {
	db *sql.DB
}

func NewCourseRepo(db *sql.DB) *courseRepo {
	return &courseRepo{db: db}
}

func (c *courseRepo) GetCourses() ([]models.Course, error) {
	var courses []models.Course

	rows, err := c.db.Query("select id, name, description from courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Course
		err := rows.Scan(&c.Id, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}

	return courses, nil
}
