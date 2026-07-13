package models

import (
	"encoding/json"
	"fmt"
)

type Lesson struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type LessonDetail struct {
	Lesson
	CourseName string `json:"course_name"`
	CourseId   int    `json:"course_id"`
	Completed  *bool  `json:"completed,omitempty"`
}

type Answer struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type LessonAnswer struct {
	AnswerId   int    `json:"answer_id"`
	Text       string `json:"text"`
	QuestionId int    `json:"-"`
}

type NewTest struct {
	Id            int      `json:"-"`
	Question      string   `json:"question"`
	AnswerOptions []Answer `json:"answer_options"`
}

type LessonTest struct {
	QuestionId    int            `json:"question_id" redis:"question_id"`
	Question      string         `json:"question" redis:"question"`
	AnswerOptions []LessonAnswer `json:"answer_options" redis:"answer_options"`
}

type AnswerOptionsView struct {
	AnswerId int    `json:"answer_id"`
	Text     string `json:"text"`
}

type AnswerOptionsList []AnswerOptionsView

func (a *AnswerOptionsList) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion fail")
	}

	return json.Unmarshal(source, a)
}

type StudentTestView struct {
	QuestionId    int               `json:"question_id"`
	Question      string            `json:"question"`
	AnswerOptions AnswerOptionsList `json:"answer_options"`
}

type AdminAnswersList []Answer

func (a *AdminAnswersList) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion fail")
	}

	return json.Unmarshal(source, a)
}

type AdminTestView struct {
	QuestionId    int              `json:"question_id"`
	Question      string           `json:"question"`
	AnswerOptions AdminAnswersList `json:"answer_options"`
}

type CorrectAnswers struct {
	QuestionId     int    `json:"question_id"`
	Question       string `json:"question"`
	CorrectChoice  int    `json:"correct_choice"`
	SelectedChoice int    `json:"selected_choice"`
	IsCorrect      bool   `json:"is_correct"`
}

type SubmitTest struct {
	QuestionId     int `json:"question_id"`
	SelectedChoice int `json:"selected_choice"`
}

type ResultsOfTest struct {
	TotalQuestions int              `json:"total_questions"`
	Point          int              `json:"points"`
	Results        []CorrectAnswers `json:"results"`
}
