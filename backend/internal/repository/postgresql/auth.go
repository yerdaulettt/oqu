package postgresql

import (
	"database/sql"
	"os"
	"time"

	"oqu/internal/models"
	"oqu/utils"

	"github.com/golang-jwt/jwt/v5"
)

type authRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *authRepo {
	return &authRepo{db: db}
}

func (r *authRepo) Register(u *models.UserRegister) (int, error) {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return -1, err
	}

	var id int
	query := `insert into users (username, password, role) values ($1, $2, $3) returning id`
	err = r.db.QueryRow(query, u.Username, hashedPassword, u.Role).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *authRepo) Login(u *models.UserLogin) (string, error) {
	var userFromDB models.UserResponseDB

	query := `select username, password, role from users where username = $1`
	err := r.db.QueryRow(query, u.Username).Scan(&userFromDB.Username, &userFromDB.PasswordHash, &userFromDB.Role)
	if err != nil {
		return "", err
	}

	if utils.VerifyPassword(userFromDB.PasswordHash, u.Password) != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userFromDB.Username,
		"role":     userFromDB.Role,
		"exp":      8 * time.Hour,
	})

	s, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return s, nil
}
