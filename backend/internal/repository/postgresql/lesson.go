package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"oqu/internal/models"
)

type lessonRepo struct {
	db *sql.DB
}

func NewLessonRepo(db *sql.DB) *lessonRepo {
	return &lessonRepo{db: db}
}

func (r *lessonRepo) GetLesson(id, userId int) (*models.LessonDetail, error) {
	query := `
	with
	active_user as
		(select id, role from users where id = $1)

	select
		l.id, l.name, l.content, c.name as course_name, c.id as course_id,
		case when (select role from active_user) = 'user' then coalesce(r.completed, false) end as completed
	from (lessons as l left join rating as r on l.id = r.lesson_id and r.user_id = (select id from active_user))
	join courses as c on l.course_id = c.id where l.id = $2`

	var l models.LessonDetail
	err := r.db.QueryRow(query, userId, id).Scan(&l.Id, &l.Name, &l.Content, &l.CourseName, &l.CourseId, &l.Completed)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (r *lessonRepo) GetComments(lessonId, userId int) ([]models.Comment, error) {
	var comments []models.Comment

	query := `
	with comment_results as (
		select c.id, c.content, c.user_id, u.name, coalesce(sum(v.voted::integer), 0) as votes, c.posted_at from
		comments as c left join comment_votes as v on c.id = v.comment_id join users as u on c.user_id = u.id
		where c.lesson_id = $1 group by c.id, u.id
	), active_user (id) as (
		values ($2::integer)
	)

	select cr.id, cr.content, cr.name, cr.votes,
	coalesce((select voted from comment_votes where user_id = (select id from active_user) and comment_id = cr.id), false) as voted, cr.posted_at
	from comment_results as cr left join comment_votes as v on cr.id = v.comment_id
	left join users as u on cr.user_id = u.id and u.id = (select id from active_user)`

	rows, err := r.db.Query(query, lessonId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.Id, &c.Content, &c.AuthorName, &c.Votes, &c.Voted, &c.PostedAt)
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

func (r *lessonRepo) GetTest(lessonId, userId int) ([]models.StudentTestView, error) {
	query := `
	with correct_answers as
		(select a.question_id, a.id from (lesson_tests as lt join questions as q on lt.id = q.test_id and lt.lesson_id = $1)
		join answers as a on q.id = a.question_id and a.is_correct = true),
	tests as
		(select q.id, q.text, json_agg(json_build_object('answer_id', a.id, 'text', a.text)) as answer_options from
		(lesson_tests as lt join questions as q on lt.id = q.test_id and lt.lesson_id = $2) join answers as a on q.id = a.question_id group by q.id)
	
	select
		t.id, t.text, ts.selected_choice, case when ts.selected_choice is not null then ca.id end as correct_choice,
		case when ts.selected_choice = ca.id then true when ts.selected_choice <> ca.id then false end as is_correct,
		t.answer_options from
	(tests as t join correct_answers as ca on t.id = ca.question_id) left join test_submits as ts on t.id = ts.question_id and ts.user_id = $3`

	var tests []models.StudentTestView

	rows, err := r.db.Query(query, lessonId, lessonId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var test models.StudentTestView
		err := rows.Scan(&test.QuestionId, &test.Question, &test.SelectedChoice, &test.CorrectChoice, &test.IsCorrect, &test.AnswerOptions)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}

	return tests, nil
}

func (r *lessonRepo) ResetTest(lessonId, userId int) error {
	t, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()

	var testId int
	err = t.QueryRow(`select id from lesson_tests where lesson_id = $1`, lessonId).Scan(&testId)
	if err != nil {
		return err
	}

	_, err = t.Exec(`delete from test_submits where test_id = $1 and user_id = $2`, testId, userId)
	if err != nil {
		return err
	}

	_, err = t.Exec(`delete from test_results where test_id = $1 and user_id = $2`, testId, userId)
	if err != nil {
		return err
	}

	if err := t.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *lessonRepo) GetCorrectAnswers(lessonId int) ([]models.CorrectAnswers, error) {
	query := `
	select q.id, a.text, a.id as correct_choice from (lesson_tests as lt join questions as q on lt.id = q.test_id and lt.lesson_id = $1)
	join answers as a on q.id = a.question_id and a.is_correct = true order by q.id`

	var correctAns []models.CorrectAnswers

	rows, err := r.db.Query(query, lessonId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ca models.CorrectAnswers
		err := rows.Scan(&ca.QuestionId, &ca.Question, &ca.CorrectChoice)
		if err != nil {
			return nil, err
		}
		correctAns = append(correctAns, ca)
	}

	return correctAns, nil
}

func (r *lessonRepo) SubmitTest(lessonId, userId int, completed bool, st []models.SubmitTest) error {
	t, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()

	var testId int
	err = t.QueryRow(`select id from lesson_tests where lesson_id = $1`, lessonId).Scan(&testId)
	if err != nil {
		return err
	}

	_, err = t.Exec(`insert into test_results (test_id, user_id, completed) values ($1, $2, $3)`, testId, userId, completed)
	if err != nil {
		return err
	}

	var query strings.Builder
	query.WriteString(`insert into test_submits (test_id, user_id, question_id, selected_choice) values`)

	params := []any{}
	cnt := 1

	for i, s := range st {
		fmt.Fprintf(&query, " ($%d, $%d, $%d, $%d)", cnt, cnt+1, cnt+2, cnt+3)
		cnt += 4

		if i != len(st)-1 {
			query.WriteString(`,`)
		}

		params = append(params, testId, userId, s.QuestionId, s.SelectedChoice)
	}

	_, err = t.Exec(query.String(), params...)
	if err != nil {
		return err
	}

	err = t.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *lessonRepo) IsTestCompleted(lessonId, userId int) bool {
	query := `
	select coalesce(tr.completed, false) from lesson_tests as lt left join
	test_results as tr on lt.id = tr.test_id where lt.lesson_id = $1 and tr.user_id = $2`

	var completed bool
	err := r.db.QueryRow(query, lessonId, userId).Scan(&completed)
	if err != nil {
		log.Println(err)
		return false
	}

	return completed
}
