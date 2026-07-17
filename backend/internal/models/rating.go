package models

type Rating struct {
	CourseId         int    `json:"course_id"`
	CourseName       string `json:"course_name"`
	TotalLessons     int    `json:"total_lessons"`
	CompletedLessons int    `json:"completed_lessons"`
	ScorePercentage  int    `json:"score_percentage"`
}
