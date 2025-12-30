package repository

import (
	"database/sql"
	"fmt"

	"github.com/CreateLab/laritmo/internal/models"
	sq "github.com/Masterminds/squirrel"
)

type LectureRepository struct {
	db *sql.DB
}

func NewLectureRepository(db *sql.DB) *LectureRepository {
	return &LectureRepository{db: db}
}


func (r *LectureRepository) GetAll(courseID *int) ([]models.Lecture, error) {
	builder := sq.Select("id", "course_id", "week", "title", "content", "github_url", "created_at", "updated_at").
		From("lectures").
		OrderBy("course_id", "week")

	if courseID != nil {
		builder = builder.Where(sq.Eq{"course_id": *courseID})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get lectures: %w", err)
	}
	defer rows.Close()

	var lectures []models.Lecture
	for rows.Next() {
		var l models.Lecture
		if err := rows.Scan(&l.ID, &l.CourseID, &l.Week, &l.Title, &l.Content, &l.GithubURL, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error for lecture: %w", err)
		}
		lectures = append(lectures, l)
	}

	return lectures, nil
}


func (r *LectureRepository) GetByID(id int) (*models.Lecture, error) {
	query, args, err := sq.Select("id", "course_id", "week", "title", "content", "github_url", "created_at", "updated_at").
		From("lectures").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var l models.Lecture
	err = r.db.QueryRow(query, args...).Scan(&l.ID, &l.CourseID, &l.Week, &l.Title, &l.Content, &l.GithubURL, &l.CreatedAt, &l.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get lecture: %w", err)
	}

	return &l, nil
}

func (r *LectureRepository) Create(courseID, week int, title, content, githubURL string) (*models.Lecture, error) {
	query, args, _ := sq.Insert("lectures").
		Columns("course_id", "week", "title", "content", "github_url").
		Values(courseID, week, title, content, githubURL).
		ToSql()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	return &models.Lecture{
		ID:        int(id),
		CourseID:  courseID,
		Week:      week,
		Title:     title,
		Content:   content,
		GithubURL: &githubURL,
	}, nil
}

func (r *LectureRepository) Update(id, courseID, week int, title, content, githubURL string) error {
	query, args, _ := sq.Update("lectures").
		Set("course_id", courseID).
		Set("week", week).
		Set("title", title).
		Set("content", content).
		Set("github_url", githubURL).
		Where(sq.Eq{"id": id}).
		ToSql()

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *LectureRepository) Delete(id int) error {
	query, args, _ := sq.Delete("lectures").
		Where(sq.Eq{"id": id}).
		ToSql()

	_, err := r.db.Exec(query, args...)
	return err
}
