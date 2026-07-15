package models

type User struct {
	Id       int    `json:"id" redis:"id"`
	Name     string `json:"name" redis:"name"`
	Username string `json:"username" redis:"username"`
	Role     string `json:"role" redis:"role"`
}

type UserRegister struct {
	Name     string `json:"name"`
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

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type RefreshRequest struct {
	Refresh string `json:"refresh"`
}
