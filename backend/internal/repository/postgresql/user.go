package postgresql

import (
	"database/sql"
	"errors"

	"oqu/internal/models"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetProfileInfo(userId int) (*models.User, error) {
	var profile models.User

	query := `select id, name, username, role from users where id = $1`
	err := r.db.QueryRow(query, userId).Scan(&profile.Id, &profile.Name, &profile.Username, &profile.Role)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *userRepo) UpdateProfile(params []any, columns []string) (*models.User, error) {
	query := updateQuery("users", "id", columns)
	query.WriteString(" returning id, name, username, role")

	var u models.User
	err := r.db.QueryRow(query.String(), params...).Scan(&u.Id, &u.Name, &u.Username, &u.Role)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) GetMyClasses(userId int) ([]models.Course, error) {
	var courses []models.Course

	query := `select c.id, c.name, c.description from courses as c join enrollments as e on c.id = e.course_id where user_id = $1`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Course
		err := rows.Scan(&c.Id, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}

	return courses, nil
}

func (r *userRepo) GetAllCoursesRating(userId int) ([]models.Rating, error) {
	var ratings []models.Rating

	query := `
	with user_ratings as (
		select e.course_id, coalesce(sum(r.completed::integer), 0) as completed_lessons
		from enrollments as e left join rating as r on e.user_id = r.user_id and e.course_id = r.course_id where e.user_id = $1 group by e.course_id
	)
	
	select c.id, c.name, (select count(*) from lessons where course_id = c.id) as total_lessons, ur.completed_lessons from
	user_ratings as ur join courses as c on ur.course_id = c.id`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rating models.Rating
		err := rows.Scan(&rating.CourseId, &rating.CourseName, &rating.TotalLessons, &rating.CompletedLessons)
		if err != nil {
			return nil, err
		}

		if rating.TotalLessons != 0 {
			rating.ScorePercentage = (rating.CompletedLessons * 100 / rating.TotalLessons)
		}

		ratings = append(ratings, rating)
	}

	return ratings, nil
}

func (r *userRepo) UsernameExists(username string) (bool, error) {
	var id int

	err := r.db.QueryRow(`select id from users where username = $1`, username).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if id == 0 {
		return false, nil
	}

	return true, nil
}
