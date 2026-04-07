package course

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

func (c *courseRepo) GetCourses() ([]models.Course, error) {
	var courses []models.Course

	rows, err := c.db.Query("select id, name, description from courses")
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

func (c *courseRepo) GetCourseById(id int) (*models.Course, error) {
	var course models.Course

	query := `select id, name, description from courses where id = $1`
	err := c.db.QueryRow(query, id).Scan(&course.Id, &course.Name, &course.Description)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (c *courseRepo) GetCourseLessons(id int) ([]models.Lesson, error) {
	var courseLessons []models.Lesson

	query := `select l.id, l.name, l.content from courses as c join lessons as l on c.id = l.course_id where c.id = $1`
	rows, err := c.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var lesson models.Lesson
		err := rows.Scan(&lesson.Id, &lesson.Name, &lesson.Content)
		if err != nil {
			return nil, err
		}
		courseLessons = append(courseLessons, lesson)
	}

	return courseLessons, nil
}

func (c *courseRepo) DeleteCourse(id int) (*models.Course, error) {
	var deleted models.Course
	query := `delete from courses where id = $1 returning *`
	err := c.db.QueryRow(query, id).Scan(&deleted.Id, &deleted.Name, &deleted.Description)
	if err != nil {
		return nil, err
	}

	return &deleted, nil
}
