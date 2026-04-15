package postgresql

import (
	"database/sql"
	"fmt"
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

func (r *userRepo) getCourseDetails(id int) (string, int, error) {
	var name string
	err := r.db.QueryRow("select name from courses where id = $1", id).Scan(&name)
	if err != nil {
		return "", 0, err
	}

	var totalLessons int
	err = r.db.QueryRow("select count(*) from lessons where course_id = $1", id).Scan(&totalLessons)
	if err != nil {
		return "", 0, err
	}

	return name, totalLessons, nil
}

func (r *userRepo) GetAllCoursesRating(userId int) ([]models.Rating, error) {
	var ratings []models.Rating

	query := `select e.course_id, sum(r.lesson_score) as total from
	enrollments as e join rating as r on e.course_id = r.course_id where e.user_id = $1 group by e.course_id;`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var courseId, sum int
		err := rows.Scan(&courseId, &sum)
		if err != nil {
			return nil, err
		}

		courseName, totalLessons, err := r.getCourseDetails(courseId)
		if err != nil {
			return nil, err
		}

		totalScore := fmt.Sprintf("Completed %d lessons from %d totally. Your score is %f", sum, totalLessons, float64(sum)/float64(totalLessons)*100)
		ratings = append(ratings, models.Rating{CourseName: courseName, TotalScore: totalScore})
	}

	return ratings, nil
}
