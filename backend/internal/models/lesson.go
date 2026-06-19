package models

type Lesson struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type LessonDetail struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	CourseName string `json:"course_name"`
	CourseId   int    `json:"course_id"`
}
