package models

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponseDB struct {
	Id           int
	Username     string
	PasswordHash string
	Role         string
}
