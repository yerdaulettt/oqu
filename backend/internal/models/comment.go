package models

type Comment struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

type ModeratorCommentView struct {
	Id         int    `json:"id"`
	Content    string `json:"content"`
	CourseName string `json:"course_name"`
	LessonName string `json:"lesson_name"`
}
