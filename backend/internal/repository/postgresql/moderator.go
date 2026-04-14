package postgresql

import (
	"database/sql"

	"oqu/internal/models"
)

type moderatorRepo struct {
	db *sql.DB
}

func NewModeratorRepo(db *sql.DB) *moderatorRepo {
	return &moderatorRepo{db: db}
}

func (r *moderatorRepo) ViewComments() ([]models.ModeratorCommentView, error) {
	var comments []models.ModeratorCommentView

	query := `select c.id, c.content as comment, co.name as course, l.name as lesson from
	(comments as c join lessons as l on c.lesson_id = l.id) join courses as co on course_id = co.id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.ModeratorCommentView
		err := rows.Scan(&c.Id, &c.Content, &c.CourseName, &c.LessonName)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *moderatorRepo) DeleteComment(id int) (*models.Comment, error) {
	var deleted models.Comment

	query := `delete from comments where id = $1 returning id, content`
	err := r.db.QueryRow(query, id).Scan(&deleted.Id, &deleted.Content)
	if err != nil {
		return nil, err
	}

	return &deleted, nil
}
