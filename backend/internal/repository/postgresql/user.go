package postgresql

import (
	"database/sql"
	"oqu/internal/models"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetProfileInfo(userId int) (*models.User, error) {
	var profile models.User

	query := `select id, name, username, role from users where id = $1`
	err := r.db.QueryRow(query, userId).Scan(&profile.Id, &profile.Name, &profile.Username, &profile.Role)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *userRepo) GetMyClasses(userId int) ([]models.Course, error) {
	var courses []models.Course

	query := `select c.id, c.name, c.description from courses as c join enrollments as e on c.id = e.course_id where user_id = $1`
	rows, err := r.db.Query(query, userId)
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
