package service

import "oqu/internal/models"

type CourseService interface {
	Get() []models.Course
	GetById(id int) *models.Course
	GetCourseLessons(id int) []models.Lesson
	EnrollInClass(classId int, userId int) error
}

type LessonService interface {
	GetComments(id int) []models.Comment
	PostComment(lessonId int, c *models.Comment) bool
}

type AuthService interface {
	Register(u *models.UserRegister) int
	Login(u *models.UserLogin) string
}

type AdminService interface {
	GetUsers() []models.User
	MakeCourse(c *models.Course) int
	Delete(id int) *models.Course
}

type ModeratorService interface {
	ViewComments() ([]models.ModeratorCommentView, error)
	DeleteComment(id int) (*models.Comment, error)
}

type UserService interface {
	GetProfileInfo(userId int) (*models.User, error)
	GetMyClasses(userId int) ([]models.Course, error)
}
