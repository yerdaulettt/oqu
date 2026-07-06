package service

import "oqu/internal/models"

type CourseService interface {
	Get() ([]models.Course, error)
	GetById(id int) (*models.Course, error)
	GetCourseLessons(id int) ([]models.Lesson, error)
	EnrollInClass(classId int, userId int) error
	ResetRating(courseId, userId int) error
}

type LessonService interface {
	GetLesson(id int) (*models.LessonDetail, error)
	GetComments(lessonId, userId int) ([]models.Comment, error)
	PostComment(lessonId int, userId int, c *models.Comment) bool
	Score(lessonId, userId int) error
	ResetScore(lessonId, userId int) error
	GetTest(lessonId int) ([]models.StudentTestView, error)
	SubmitTest(lessonId int, st []models.SubmitTest) (*models.ResultsOfTest, error)
}

type CommentService interface {
	Vote(userId, commentId int) error
	ModifyVote(userId, commentId int) error
}

type AuthService interface {
	Register(u *models.UserRegister) int
	Login(u *models.UserLogin) string
}

type AdminService interface {
	GetUsers() []models.User
	MakeCourse(c *models.Course) int
	Delete(id int) *models.Course
	AddLesson(courseId int, l *models.Lesson) (int, error)
	AddTest(lessonId int, t []*models.NewTest) error
	GetTest(lessonId int) ([]models.AdminTestView, error)
}

type ModeratorService interface {
	ViewComments() ([]models.ModeratorCommentView, error)
	DeleteComment(id int) (*models.Comment, error)
}

type UserService interface {
	GetProfileInfo(userId int) (*models.User, error)
	GetMyClasses(userId int) ([]models.Course, error)
	GetAllCoursesRating(userId int) ([]models.Rating, error)
}
