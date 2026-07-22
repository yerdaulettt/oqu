package service

import (
	"database/sql"
	"errors"
	"log"
	"slices"
	"sort"

	"oqu/internal/models"
	"oqu/internal/repository"
)

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(r repository.AdminRepository) *adminService {
	return &adminService{repo: r}
}

func (s *adminService) GetUsers() []models.User {
	users, err := s.repo.GetUsers()
	if err != nil {
		log.Println(err)
		return nil
	}
	return users
}

func (s *adminService) DeleteUser(userId int) (*models.User, error) {
	user, err := s.repo.DeleteUser(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErr
		}

		log.Println(err)
		return nil, internalErr
	}

	return user, nil
}

func (s *adminService) UpdateUserRole(userId int, role string) (*models.User, error) {
	if !slices.Contains([]string{"user", "admin", "moderator"}, role) {
		return nil, IncorrectRole
	}

	user, err := s.repo.UpdateUserRole(userId, role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErr
		}

		log.Println(err)
		return nil, internalErr
	}

	return user, nil
}

func (s *adminService) MakeCourse(c *models.NewCourse) int {
	id, err := s.repo.MakeCourse(c)
	if err != nil {
		log.Println(err)
		return 0
	}

	return id
}

func (s *adminService) UpdateCourse(c *models.NewCourse, courseId int) (*models.Course, error) {
	params := []any{}
	columns := []string{}

	if c.Name != "" {
		columns = append(columns, "name")
		params = append(params, c.Name)
	}

	if c.Description != "" {
		columns = append(columns, "description")
		params = append(params, c.Description)
	}

	params = append(params, courseId)

	course, err := s.repo.UpdateCourse(params, columns)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return course, nil
}

func (s *adminService) Delete(id int) *models.Course {
	result, err := s.repo.DeleteCourse(id)
	if err != nil {
		log.Println("course service error:", err)
		return nil
	}

	return result
}

func (s *adminService) AddLesson(courseId int, l *models.NewLesson) (int, error) {
	id, err := s.repo.AddLesson(courseId, l)

	if err != nil {
		log.Println("Admin service AddLesson():", err)
		return id, internalErr
	}

	return id, nil
}

func (s *adminService) UpdateLesson(lessonId int, l *models.NewLesson) (*models.Lesson, error) {
	params := []any{}
	columns := []string{}

	if l.Name != "" {
		params = append(params, l.Name)
		columns = append(columns, "name")
	}

	if l.Content != "" {
		params = append(params, l.Content)
		columns = append(columns, "content")
	}

	params = append(params, lessonId)

	lesson, err := s.repo.UpdateLesson(params, columns)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErr
		}

		log.Println(err)
		return nil, internalErr
	}

	return lesson, nil
}

func (s *adminService) DeleteLesson(lessonId int) (*models.Lesson, error) {
	lesson, err := s.repo.DeleteLesson(lessonId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErr
		}

		log.Println(err)
		return nil, internalErr
	}

	return lesson, nil
}

func (s *adminService) AddTest(lessonId int, t []*models.NewTest) error {
	err := s.repo.AddTest(lessonId, t)
	if err != nil {
		log.Println(err)
		return internalErr
	}

	return nil
}

func (s *adminService) GetTest(lessonId int) ([]models.AdminTestView, error) {
	tests, err := s.repo.GetTest(lessonId)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return tests, nil
}

func (s *adminService) UpdateTest(lessonId int, newTest []models.AdminTestView) ([]models.AdminTestView, error) {
	testFromDB, err := s.repo.GetTest(lessonId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundErr
		}

		log.Println(err)
		return nil, internalErr
	}

	sort.Slice(newTest, func(i, j int) bool {
		sort.Slice(newTest[i].AnswerOptions, func(i1, i2 int) bool {
			return newTest[i].AnswerOptions[i1].AnswerId < newTest[i].AnswerOptions[i2].AnswerId
		})

		return newTest[i].QuestionId < newTest[j].QuestionId
	})

	var cnt int
	questionParams := []any{}
	answerParams := []any{}

	for i, t := range newTest {
		if t.QuestionId == testFromDB[i].QuestionId {
			cnt = i
		} else {
			for k := i + 1; k < len(testFromDB); k++ {
				if t.QuestionId == testFromDB[k].QuestionId {
					cnt = k
					break
				}
			}
		}

		if t.Question != "" && t.Question != testFromDB[cnt].Question {
			questionParams = append(questionParams, t.Question, t.QuestionId)
			testFromDB[cnt].Question = t.Question
		}

		for j, a := range t.AnswerOptions {
			if a.AnswerId != testFromDB[cnt].AnswerOptions[j].AnswerId {
				continue
			}

			var answerHelper []any

			if a.Text != "" && a.Text != testFromDB[cnt].AnswerOptions[j].Text {
				answerHelper = append(answerHelper, a.Text)
				testFromDB[cnt].AnswerOptions[j].Text = a.Text
			} else {
				answerHelper = append(answerHelper, nil)
			}

			if a.IsCorrect != testFromDB[cnt].AnswerOptions[j].IsCorrect {
				answerHelper = append(answerHelper, a.IsCorrect)
				testFromDB[cnt].AnswerOptions[j].IsCorrect = a.IsCorrect
			} else {
				answerHelper = append(answerHelper, nil)
			}

			if answerHelper[0] == nil && answerHelper[1] == nil {
				continue
			} else {
				answerHelper = append(answerHelper, a.AnswerId)
				answerParams = append(answerParams, answerHelper...)
			}
		}
	}

	if len(questionParams) == 0 && len(answerParams) == 0 {
		return nil, UpdateErr
	}

	err = s.repo.UpdateTest(questionParams, answerParams)
	if err != nil {
		log.Println(err)
		return nil, internalErr
	}

	return testFromDB, nil
}

func (s *adminService) DeleteTest(lessonId int) error {
	err := s.repo.DeleteTest(lessonId)
	if err != nil {
		log.Println(err)
		return internalErr
	}

	return nil
}
