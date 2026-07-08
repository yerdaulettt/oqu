package models

import (
	"encoding/json"
	"fmt"
)

type Course struct {
	Id          int    `json:"id" redis:"id"`
	Name        string `json:"name" redis:"name"`
	Description string `json:"description" redis:"description"`
}

type CourseLesson struct {
	Id        int    `json:"lesson_id"`
	Name      string `json:"lesson_name"`
	Completed *bool  `json:"completed,omitempty"`
}

type CourseLessonList []CourseLesson

func (cl *CourseLessonList) Scan(src any) error {
	if src == nil {
		*cl = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion fail")
	}

	return json.Unmarshal(source, cl)
}

type CourseDetails struct {
	Course
	Lessons CourseLessonList `json:"lessons"`
}
