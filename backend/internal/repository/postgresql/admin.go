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

func (r *adminRepo) MakeCourse(c *models.Course) (int, error) {
	var id int
	query := `insert into courses values(default, $1, $2) returning id`
	err := r.db.QueryRow(query, c.Name, c.Description).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
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

func (r *adminRepo) AddLesson(courseId int, l *models.Lesson) (int, error) {
	var id int

	query := `insert into lessons (name, content, course_id) values ($1, $2, $3) returning id`
	err := r.db.QueryRow(query, l.Name, l.Content, courseId).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
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

func (r *adminRepo) AddTest(lessonId int, t []*models.NewTest) error {
	query := `insert into questions (text, lesson_id) values `
	cnt := 1
	params := []any{}

	for i, q := range t {
		query += `($` + strconv.Itoa(cnt) + `, $` + strconv.Itoa(cnt+1) + `)`
		cnt += 2
		params = append(params, q.Question, lessonId)

		if i+1 == len(t) {
			query += ` returning id;`
		} else {
			query += `, `
		}
	}

	rows, err := r.db.Query(query, params...)
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
		t[cnt].Id = id
		cnt += 1
	}

	query, params = addTestHelper(t)

	_, err = r.db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r *adminRepo) GetTest(lessonId int) ([]models.AdminTestView, error) {
	query := `select q.id, q.text, json_agg(json_build_object('text', a.text, 'is_correct', a.is_correct)) as answer_options from
	questions as q join answers as a on q.id = a.question_id where q.lesson_id = $1 group by q.id`

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
