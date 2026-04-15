package repository

import "oqu/internal/models"

type CourseRepository interface {
	GetCourses() ([]models.Course, error)
	GetCourseById(id int) (*models.Course, error)
	GetCourseLessons(id int) ([]models.Lesson, error)
	EnrollInClass(classId int, userId int) error
}

type LessonRepository interface {
	GetComments(id int) ([]models.Comment, error)
	PostComment(lessonId int, userId int, c *models.Comment) (bool, error)
}

type CommentRepository interface {
	Vote(commentId int) error
}

type AuthRepository interface {
	Register(u *models.UserRegister) (int, error)
	Login(u *models.UserLogin) (string, error)
}

type AdminRepository interface {
	GetUsers() ([]models.User, error)
	MakeCourse(c *models.Course) (int, error)
	DeleteCourse(id int) (*models.Course, error)
}

type ModeratorRepository interface {
	ViewComments() ([]models.ModeratorCommentView, error)
	DeleteComment(id int) (*models.Comment, error)
}

type UserRepository interface {
	GetProfileInfo(userId int) (*models.User, error)
	GetMyClasses(userId int) ([]models.Course, error)
}
