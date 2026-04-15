package postgresql

import "database/sql"

type commentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *commentRepo {
	return &commentRepo{db: db}
}

func (r *commentRepo) Vote(commentId int) error {
	query := `update comments set votes = votes + 1 where id = $1`
	_, err := r.db.Exec(query, commentId)
	if err != nil {
		return err
	}

	return nil
}
