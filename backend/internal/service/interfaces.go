package service

import "oqu/internal/models"

type CourseService interface {
	Get() []models.Course
	GetById(id int) *models.Course
	GetCourseLessons(id int) []models.Lesson
	MakeCourse(c *models.Course) int
	Delete(id int) *models.Course
}

type LessonService interface {
	GetComments(id int) []models.Comment
}

type AuthService interface {
	Register(u *models.UserRegister) int
	Login(u *models.UserLogin) string
}

type AdminService interface {
	GetUsers() []models.User
}
