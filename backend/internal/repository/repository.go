package repository

import "oqu/internal/models"

type CourseRepository interface {
	GetCourses() ([]models.Course, error)
	GetCourseById(id int) (*models.Course, error)
	GetCourseLessons(id int) ([]models.Lesson, error)
	MakeCourse(c *models.Course) (int, error)
	DeleteCourse(id int) (*models.Course, error)
}
