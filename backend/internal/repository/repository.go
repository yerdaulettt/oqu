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
	ResetRating(courseId, userId int) error
}

type LessonRepository interface {
	GetLesson(id, userId int) (*models.LessonDetail, error)
	GetComments(lessonId, userId int) ([]models.Comment, error)
	PostComment(lessonId int, userId int, c *models.Comment) (bool, error)
	Score(lessonId, userId int) error
	ResetScore(lessonId, userId int) error
	GetTest(lessonId int) ([]models.StudentTestView, error)
	GetCorrectAnswers(lessonId int) ([]models.CorrectAnswers, error)
	SubmitTest(lessonId, userId int, completed bool, st []models.SubmitTest) error
	IsTestCompleted(lessonId, userId int) bool
}

type CommentRepository interface {
	Vote(userId, commentId int) error
	ModifyVote(userId, commentId int) error
}

type AuthRepository interface {
	Register(u *models.UserRegister) (int, error)
	GetUser(username string) (*models.UserResponseDB, error)
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

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}
