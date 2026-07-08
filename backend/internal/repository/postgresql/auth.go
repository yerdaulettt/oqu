package postgresql

import (
	"database/sql"

	"oqu/internal/models"
)

type authRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *authRepo {
	return &authRepo{db: db}
}

func (r *authRepo) Register(u *models.UserRegister) (int, error) {
	var id int
	query := `insert into users (name, username, password, role) values ($1, $2, $3, $4) returning id`
	err := r.db.QueryRow(query, u.Name, u.Username, u.Password, u.Role).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *authRepo) GetUser(username string) (*models.UserResponseDB, error) {
	var u models.UserResponseDB
	query := `select id, username, password, role from users where username = $1`

	err := r.db.QueryRow(query, username).Scan(&u.Id, &u.Username, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
