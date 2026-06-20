package postgresql

import (
	"database/sql"
	"log"

	"oqu/internal/models"
)

type lessonRepo struct {
	db *sql.DB
}

func NewLessonRepo(db *sql.DB) *lessonRepo {
	return &lessonRepo{db: db}
}

func (r *lessonRepo) GetLesson(id int) (*models.LessonDetail, error) {
	query := `select l.id, l.name, l.content, c.name, c.id from lessons as l join courses as c on l.course_id = c.id where l.id = $1`

	var l models.LessonDetail
	err := r.db.QueryRow(query, id).Scan(&l.Id, &l.Name, &l.Content, &l.CourseName, &l.CourseId)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (r *lessonRepo) GetComments(id int) ([]models.Comment, error) {
	var comments []models.Comment

	query := `select c.id, c.content, sum(voted::integer) as votes from
	comments as c join comment_votes as v on c.id = v.comment_id where lesson_id = $1 group by c.id`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.Id, &c.Content, &c.Votes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *lessonRepo) PostComment(lessonId int, userId int, c *models.Comment) (bool, error) {
	query := `insert into comments (content, lesson_id, user_id) values ($1, $2, $3)`

	if res, err := r.db.Exec(query, c.Content, lessonId, userId); err != nil {
		log.Println(res)
		return false, err
	}

	return true, nil
}

func (r *lessonRepo) Score(lessonId, userId int) error {
	query := `insert into rating (lesson_id, user_id, completed, course_id) values ($1, $2, $3, (select course_id from lessons where id = $4))`

	_, err := r.db.Exec(query, lessonId, userId, true, lessonId)
	if err != nil {
		return err
	}

	return nil
}

func (r *lessonRepo) ResetScore(lessonId, userId int) error {
	query := `delete from rating where lesson_id = $1 and user_id = $2`

	_, err := r.db.Exec(query, lessonId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *lessonRepo) GetTest(lessonId int) ([]models.StudentTestView, error) {
	query := `select q.id, q.text, json_agg(json_build_object('answer_id', a.id, 'text', a.text)) as answer_options from
	questions as q join answers as a on q.id = a.question_id where q.lesson_id = $1 group by q.id`
	var tests []models.StudentTestView

	rows, err := r.db.Query(query, lessonId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var test models.StudentTestView
		err := rows.Scan(&test.QuestionId, &test.Question, &test.AnswerOptions)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}

	return tests, nil
}

func (r *lessonRepo) GetCorrectAnswers(lessonId int) ([]models.CorrectAnswers, error) {
	query := `select q.id, a.id from questions as q join answers as a on q.id = a.question_id where lesson_id = $1 and is_correct = true`
	var correctAns []models.CorrectAnswers

	rows, err := r.db.Query(query, lessonId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ca models.CorrectAnswers
		err := rows.Scan(&ca.QuestionId, &ca.CorrectOptionId)
		if err != nil {
			return nil, err
		}
		correctAns = append(correctAns, ca)
	}

	return correctAns, nil
}
