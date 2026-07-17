package postgresql

import (
	"database/sql"

	"oqu/internal/models"
)

type courseRepo struct {
	db *sql.DB
}

func NewCourseRepo(db *sql.DB) *courseRepo {
	return &courseRepo{db: db}
}

func (r *courseRepo) GetCourses() ([]models.Course, error) {
	var courses []models.Course

	rows, err := r.db.Query("select id, name, description from courses")
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

func (r *courseRepo) GetCourseById(id, userId int) (*models.CourseDetails, error) {
	var c models.CourseDetails
	query := `
	with
	active_user as
		(select id, role from users where id = $1),
	course_lessons as
		(select
			l.id as lesson_id, l.name as lesson_name,
			case when (select role from active_user) = 'user' then coalesce(r.completed, false) end as completed
		from lessons as l left join rating as r on l.id = r.lesson_id and r.user_id = (select id from active_user) where l.course_id = $2)

	select id, name, description, (select json_agg(to_jsonb(course_lessons)) as lessons from course_lessons) from courses where id = $3`

	err := r.db.QueryRow(query, userId, id, id).Scan(&c.Id, &c.Name, &c.Description, &c.Lessons)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *courseRepo) EnrollInClass(classId int, userId int) error {
	query := `insert into enrollments (course_id, user_id) values ($1, $2) on conflict do nothing`
	_, err := r.db.Exec(query, classId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *courseRepo) Unenroll(courseId, userId int) error {
	query := `delete from enrollments where course_id = $1 and user_id = $2`

	_, err := r.db.Exec(query, courseId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *courseRepo) ResetRating(courseId, userId int) error {
	query := `delete from rating where course_id = $1 and user_id = $2`

	_, err := r.db.Exec(query, courseId, userId)
	if err != nil {
		return err
	}

	return nil
}
