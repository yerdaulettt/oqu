package postgresql

import (
	"database/sql"
	"strconv"

	"oqu/internal/models"
)

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) GetUsers() ([]models.User, error) {
	var users []models.User

	rows, err := r.db.Query("select id, name, username, role from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.Id, &u.Name, &u.Username, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *adminRepo) DeleteUser(userId int) (*models.User, error) {
	query := `delete from users where id = $1 returning id, name, username, role`

	var u models.User
	if err := r.db.QueryRow(query, userId).Scan(&u.Id, &u.Name, &u.Username, &u.Role); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *adminRepo) UpdateUserRole(userId int, role string) (*models.User, error) {
	query := `update users set role = $1 where id = $2 returning id, name, username, role`

	var u models.User
	if err := r.db.QueryRow(query, role, userId).Scan(&u.Id, &u.Name, &u.Username, &u.Role); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *adminRepo) MakeCourse(c *models.NewCourse) (int, error) {
	var id int
	query := `insert into courses values(default, $1, $2) returning id`
	err := r.db.QueryRow(query, c.Name, c.Description).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *adminRepo) UpdateCourse(params []any, columns []string) (*models.Course, error) {
	query := updateQuery("courses", "id", columns)
	query.WriteString(" returning id, name, description")

	var c models.Course
	if err := r.db.QueryRow(query.String(), params...).Scan(&c.Id, &c.Name, &c.Description); err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *adminRepo) DeleteCourse(id int) (*models.Course, error) {
	var deleted models.Course

	query := `delete from courses where id = $1 returning *`
	err := r.db.QueryRow(query, id).Scan(&deleted.Id, &deleted.Name, &deleted.Description)
	if err != nil {
		return nil, err
	}

	return &deleted, nil
}

func (r *adminRepo) AddLesson(courseId int, l *models.NewLesson) (int, error) {
	var id int

	query := `insert into lessons (name, content, course_id) values ($1, $2, $3) returning id`
	err := r.db.QueryRow(query, l.Name, l.Content, courseId).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *adminRepo) UpdateLesson(params []any, columns []string) (*models.Lesson, error) {
	query := updateQuery("lessons", "id", columns)
	query.WriteString(" returning id, name, content")

	var l models.Lesson
	if err := r.db.QueryRow(query.String(), params...).Scan(&l.Id, &l.Name, &l.Content); err != nil {
		return nil, err
	}

	return &l, nil
}

func (r *adminRepo) DeleteLesson(lessonId int) (*models.Lesson, error) {
	query := `delete from lessons where id = $1 returning id, name, content`

	var l models.Lesson
	err := r.db.QueryRow(query, lessonId).Scan(&l.Id, &l.Name, &l.Content)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func addTestHelper(tests []*models.NewTest) (string, []any) {
	query := `insert into answers (text, is_correct, question_id) values `
	cnt := 1
	params := []any{}

	for i, t := range tests {
		for j, a := range t.AnswerOptions {
			query += `($` + strconv.Itoa(cnt) + `, $` + strconv.Itoa(cnt+1) + `, $` + strconv.Itoa(cnt+2) + `)`
			cnt += 3
			params = append(params, a.Text, a.IsCorrect, t.Id)

			if i+1 == len(tests) && j+1 == len(t.AnswerOptions) {
				query += `;`
			} else {
				query += `, `
			}
		}
	}

	return query, params
}

func (r *adminRepo) AddTest(lessonId int, nt []*models.NewTest) error {
	t, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()

	var testId int
	err = t.QueryRow(`insert into lesson_tests (lesson_id) values ($1) returning id`, lessonId).Scan(&testId)
	if err != nil {
		return err
	}

	query := `insert into questions (text, test_id) values `
	cnt := 1
	params := []any{}

	for i, q := range nt {
		query += `($` + strconv.Itoa(cnt) + `, $` + strconv.Itoa(cnt+1) + `)`
		cnt += 2
		params = append(params, q.Question, testId)

		if i+1 == len(nt) {
			query += ` returning id;`
		} else {
			query += `, `
		}
	}

	rows, err := t.Query(query, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	cnt = 0
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		nt[cnt].Id = id
		cnt += 1
	}

	query, params = addTestHelper(nt)

	_, err = t.Exec(query, params...)
	if err != nil {
		return err
	}

	if err := t.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *adminRepo) GetTest(lessonId int) ([]models.AdminTestView, error) {
	query := `
	select q.id, q.text, json_agg(json_build_object('text', a.text, 'is_correct', a.is_correct)) as answer_options
	from (lesson_tests as lt join questions as q on lt.id = q.test_id and lt.lesson_id = $1)
	join answers as a on q.id = a.question_id group by q.id`

	var adminTests []models.AdminTestView
	rows, err := r.db.Query(query, lessonId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var adminTest models.AdminTestView
		err := rows.Scan(&adminTest.QuestionId, &adminTest.Question, &adminTest.AnswerOptions)
		if err != nil {
			return nil, err
		}
		adminTests = append(adminTests, adminTest)
	}

	return adminTests, nil
}

func (r *adminRepo) DeleteTest(lessonId int) error {
	_, err := r.db.Exec(`delete from lesson_tests where lesson_id = $1`, lessonId)
	if err != nil {
		return err
	}

	return nil
}
