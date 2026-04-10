package repository

import "oqu/internal/models"

type CourseRepository interface {
	GetCourses() ([]models.Course, error)
	GetCourseById(id int) (*models.Course, error)
	GetCourseLessons(id int) ([]models.Lesson, error)
	MakeCourse(c *models.Course) (int, error)
	DeleteCourse(id int) (*models.Course, error)
}

type LessonRepository interface {
	GetComments(id int) ([]models.Comment, error)
}

type AuthRepository interface {
	Register(u *models.UserRegister) (int, error)
	Login(u *models.UserLogin) (string, error)
}

type AdminRepository interface {
	GetUsers() ([]models.User, error)
}
