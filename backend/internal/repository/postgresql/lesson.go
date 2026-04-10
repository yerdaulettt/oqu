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

	query := `select c.id, c.content from lessons as l join comments as c on l.id = c.lesson_id where l.id = $1`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.Id, &c.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *lessonRepo) PostComment(lessonId int, c *models.Comment) (bool, error) {
	query := `insert into comments (content, lesson_id) values ($1, $2)`

	if res, err := r.db.Exec(query, c.Content, lessonId); err != nil {
		log.Println(res)
		return false, err
	}

	return true, nil
}
