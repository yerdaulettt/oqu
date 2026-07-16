package service

import "oqu/internal/models"

type CourseService interface {
	Get() ([]models.Course, error)
	GetById(id, userId int) (*models.CourseDetails, error)
	EnrollInClass(classId int, userId int) error
	ResetRating(courseId, userId int) error
}

type LessonService interface {
	GetLesson(id, userId int) (*models.LessonDetail, error)
	GetComments(lessonId, userId int) ([]models.Comment, error)
	PostComment(lessonId int, userId int, c *models.Comment) bool
	Score(lessonId, userId int) error
	ResetScore(lessonId, userId int) error
	GetTest(lessonId, userId int) (*models.StudentTestView, error)
	ResetTest(lessonId, userId int) error
	SubmitTest(lessonId, userId int, st []models.SubmitTest) (*models.ResultsOfTest, error)
}

type CommentService interface {
	Vote(userId, commentId int) error
	ModifyVote(userId, commentId int) error
}

type AuthService interface {
	Register(u *models.UserRegister) int
	Login(u *models.UserLogin) (*models.Tokens, error)
	Refresh(refresh string) (string, error)
}

type AdminService interface {
	GetUsers() []models.User
	DeleteUser(userId int) (*models.User, error)
	UpdateUserRole(userId int, role string) (*models.User, error)
	MakeCourse(c *models.NewCourse) int
	UpdateCourse(c *models.NewCourse, courseId int) (*models.Course, error)
	Delete(id int) *models.Course
	AddLesson(courseId int, l *models.NewLesson) (int, error)
	UpdateLesson(lessonId int, l *models.NewLesson) (*models.Lesson, error)
	DeleteLesson(lessonId int) (*models.Lesson, error)
	AddTest(lessonId int, t []*models.NewTest) error
	GetTest(lessonId int) ([]models.AdminTestView, error)
	DeleteTest(lessonId int) error
}

type ModeratorService interface {
	ViewComments() ([]models.ModeratorCommentView, error)
	DeleteComment(id int) (*models.DeletedComment, error)
}

type UserService interface {
	GetProfileInfo(userId int) (*models.User, error)
	UpdateProfile(u *models.UserUpdate, userId int) (*models.User, error)
	GetMyClasses(userId int) ([]models.Course, error)
	GetAllCoursesRating(userId int) ([]models.Rating, error)
}
