package postgresql

import (
	"database/sql"
	"oqu/internal/models"
)

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) GetUsers() ([]models.User, error) {
	var users []models.User

	rows, err := r.db.Query("select id, name, username, role from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.Id, &u.Name, &u.Username, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *adminRepo) MakeCourse(c *models.Course) (int, error) {
	var id int
	query := `insert into courses values(default, $1, $2) returning id`
	err := r.db.QueryRow(query, c.Name, c.Description).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *adminRepo) DeleteCourse(id int) (*models.Course, error) {
	var deleted models.Course
	query := `delete from courses where id = $1 returning *`
	err := r.db.QueryRow(query, id).Scan(&deleted.Id, &deleted.Name, &deleted.Description)
	if err != nil {
		return nil, err
	}

	return &deleted, nil
}
