package postgresql

import (
	"database/sql"

	"oqu/internal/models"
)

type commentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *commentRepo {
	return &commentRepo{db: db}
}

func (r *commentRepo) GetUserId(commentId int) (int, error) {
	query := `select user_id from comments where id = $1`

	var userId int
	if err := r.db.QueryRow(query, commentId).Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}

func (r *commentRepo) UpdateComment(commentId, userId int, content string) (*models.UpdatedComment, error) {
	query := `update comments set content = $1 where id = $2 and user_id = $3 returning id, content`
	var uc models.UpdatedComment

	if err := r.db.QueryRow(query, content, commentId, userId).Scan(&uc.Id, &uc.Content); err != nil {
		return nil, err
	}

	return &uc, nil
}

func (r *commentRepo) DeleteComment(commentId, userId int) (*models.DeletedComment, error) {
	query := `delete from comments where id = $1 and user_id = $2 returning id, content`
	var c models.DeletedComment

	if err := r.db.QueryRow(query, commentId, userId).Scan(&c.Id, &c.Content); err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *commentRepo) Vote(commentId, userId int, vote bool) error {
	query := `insert into comment_votes (comment_id, user_id, voted) values ($1, $2, true) on conflict (comment_id, user_id) do update set voted = $3`

	_, err := r.db.Exec(query, commentId, userId, vote)
	if err != nil {
		return err
	}

	return nil
}
