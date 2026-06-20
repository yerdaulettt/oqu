package models

type Course struct {
	Id          int    `json:"id" redis:"id"`
	Name        string `json:"name" redis:"name"`
	Description string `json:"description" redis:"description"`
}
