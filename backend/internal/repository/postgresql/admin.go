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
