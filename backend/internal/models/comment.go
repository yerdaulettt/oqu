package models

import "time"

type Comment struct {
	Id         int       `json:"id"`
	Content    string    `json:"content"`
	AuthorName string    `json:"author_name"`
	Votes      int       `json:"votes"`
	Voted      bool      `json:"voted"`
	MyComment  bool      `json:"my_comment"`
	PostedAt   time.Time `json:"posted_at"`
}

type NewComment struct {
	Content string `json:"content"`
}

type UpdatedComment struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

type DeletedComment struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

type ModeratorCommentView struct {
	Id         int       `json:"id"`
	Content    string    `json:"content"`
	Username   string    `json:"username"`
	CourseName string    `json:"course_name"`
	CourseId   int       `json:"course_id"`
	LessonName string    `json:"lesson_name"`
	LessonId   int       `json:"lesson_id"`
	PostedAt   time.Time `json:"posted_at"`
}
