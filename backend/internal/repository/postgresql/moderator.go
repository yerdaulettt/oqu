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

	query := `
	select c.id, c.content, u.username, courses.name as course_name, courses.id as course_id,
	l.name as lesson_name, l.id as lesson_id, c.posted_at
	from (comments as c join lessons as l on c.lesson_id = l.id)
	join courses on l.course_id = courses.id join users as u on c.user_id = u.id order by posted_at desc`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.ModeratorCommentView
		err := rows.Scan(&c.Id, &c.Content, &c.Username, &c.CourseName, &c.CourseId, &c.LessonName, &c.LessonId, &c.PostedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *moderatorRepo) DeleteComment(id int) (*models.DeletedComment, error) {
	var deleted models.DeletedComment

	query := `delete from comments where id = $1 returning id, content`
	err := r.db.QueryRow(query, id).Scan(&deleted.Id, &deleted.Content)
	if err != nil {
		return nil, err
	}

	return &deleted, nil
}
