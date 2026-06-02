package postgresql

import "database/sql"

type commentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *commentRepo {
	return &commentRepo{db: db}
}

func (r *commentRepo) Vote(userId, commentId int) error {
	query := `insert into comment_votes (comment_id, user_id, voted) values ($1, $2, true) on conflict do nothing`

	_, err := r.db.Exec(query, commentId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *commentRepo) ModifyVote(userId, commentId int) error {
	query := `update comment_votes set
	voted = case voted when true then false when false then true end
	where comment_id = $1 and user_id = $2`

	_, err := r.db.Exec(query, commentId, userId)
	if err != nil {
		return err
	}

	return nil
}
