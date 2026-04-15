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

func (r *lessonRepo) GetComments(id int) ([]models.Comment, error) {
	var comments []models.Comment

	query := `select c.id, c.content, u.username, c.votes from
	(lessons as l join comments as c on l.id = c.lesson_id)
	join users as u on c.user_id = u.id where l.id = $1`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.Id, &c.Content, &c.Username, &c.Votes)
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

func (r *lessonRepo) Score(lessonId, score, userId int) error {
	var courseId int
	err := r.db.QueryRow("select course_id from lessons where id = $1", lessonId).Scan(&courseId)
	if err != nil {
		return err
	}

	query := `insert into rating (course_id, lesson_id, lesson_score, user_id) values ($1, $2, $3, $4)`
	_, err = r.db.Exec(query, courseId, lessonId, score, userId)
	if err != nil {
		return err
	}

	return nil
}
