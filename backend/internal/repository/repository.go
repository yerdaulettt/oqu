package repository

import "oqu/internal/models"

type CourseRepository interface {
	GetCourses() ([]models.Course, error)
	GetCourseById(id int) (*models.Course, error)
	GetCourseLessons(id int) ([]models.Lesson, error)
	EnrollInClass(classId int, userId int) error
	ResetRating(courseId, userId int) error
}

type LessonRepository interface {
	GetLesson(id int) (*models.LessonDetail, error)
	GetComments(id int) ([]models.Comment, error)
	PostComment(lessonId int, userId int, c *models.Comment) (bool, error)
	Score(lessonId, userId int) error
	ResetScore(lessonId, userId int) error
	GetTest(lessonId int) ([]models.StudentTestView, error)
	GetCorrectAnswers(lessonId int) ([]models.CorrectAnswers, error)
}

type CommentRepository interface {
	Vote(userId, commentId int) error
	ModifyVote(userId, commentId int) error
}

type AuthRepository interface {
	Register(u *models.UserRegister) (int, error)
	Login(u *models.UserLogin) (string, error)
}

type AdminRepository interface {
	GetUsers() ([]models.User, error)
	MakeCourse(c *models.Course) (int, error)
	DeleteCourse(id int) (*models.Course, error)
	AddLesson(courseId int, l *models.Lesson) (int, error)
	AddTest(lessonId int, t []*models.NewTest) error
	GetTest(lessonId int) ([]models.AdminTestView, error)
}

type ModeratorRepository interface {
	ViewComments() ([]models.ModeratorCommentView, error)
	DeleteComment(id int) (*models.Comment, error)
}

type UserRepository interface {
	GetProfileInfo(userId int) (*models.User, error)
	GetMyClasses(userId int) ([]models.Course, error)
	GetAllCoursesRating(userId int) ([]models.Rating, error)
}
