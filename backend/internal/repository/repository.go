package repository

import (
	"context"
	"time"

	"oqu/internal/models"
)

type CourseRepository interface {
	GetCourses() ([]models.Course, error)
	GetCourseById(id, userId int) (*models.CourseDetails, error)
	EnrollInClass(classId int, userId int) error
	Unenroll(courseId, userId int) error
	ResetRating(courseId, userId int) error
}

type LessonRepository interface {
	GetLesson(id, userId int) (*models.LessonDetail, error)
	GetComments(lessonId, userId int) ([]models.Comment, error)
	PostComment(lessonId int, userId int, c *models.Comment) (bool, error)
	Score(lessonId, userId int) error
	ResetScore(lessonId, userId int) error
	GetTest(lessonId, userId int) ([]models.StudentTestQuestions, error)
	ResetTest(lessonId, userId int) error
	GetCorrectAnswers(lessonId int) ([]models.CorrectAnswers, error)
	SubmitTest(lessonId, userId int, completed bool, st []models.SubmitTest) error
	IsTestCompleted(lessonId, userId int) bool
}

type CommentRepository interface {
	GetUserId(commentId int) (int, error)
	UpdateComment(commentId, userId int, content string) (*models.UpdatedComment, error)
	DeleteComment(commentId, userId int) (*models.DeletedComment, error)
	Vote(userId, commentId int) error
	ModifyVote(userId, commentId int) error
}

type AuthRepository interface {
	Register(u *models.UserRegister) (int, error)
	GetUser(username string) (*models.UserResponseDB, error)
}

type AdminRepository interface {
	GetUsers() ([]models.User, error)
	DeleteUser(userId int) (*models.User, error)
	UpdateUserRole(userId int, role string) (*models.User, error)
	MakeCourse(c *models.NewCourse) (int, error)
	UpdateCourse(params []any, columns []string) (*models.Course, error)
	DeleteCourse(id int) (*models.Course, error)
	AddLesson(courseId int, l *models.NewLesson) (int, error)
	UpdateLesson(params []any, columns []string) (*models.Lesson, error)
	DeleteLesson(lessonId int) (*models.Lesson, error)
	AddTest(lessonId int, nt []*models.NewTest) error
	GetTest(lessonId int) ([]models.AdminTestView, error)
	DeleteTest(lessonId int) error
}

type ModeratorRepository interface {
	ViewComments() ([]models.ModeratorCommentView, error)
	DeleteComment(id int) (*models.DeletedComment, error)
}

type UserRepository interface {
	GetProfileInfo(userId int) (*models.User, error)
	UpdateProfile(params []any, columns []string) (*models.User, error)
	UsernameExists(username string) (bool, error)
	GetMyClasses(userId int) ([]models.Course, error)
	GetAllCoursesRating(userId int) ([]models.Rating, error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}
