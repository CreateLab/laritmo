package repository

import (
	"database/sql"
	"fmt"

	"github.com/CreateLab/laritmo/internal/models"
	sq "github.com/Masterminds/squirrel"
)

type LabRepository struct {
	db *sql.DB
}

func NewLabRepository(db *sql.DB) *LabRepository {
	return &LabRepository{db: db}
}


func (r *LabRepository) GetAll(courseID *int) ([]models.Lab, error) {
	builder := sq.Select("id", "course_id", "number", "title", "description", "deadline", "max_score", "github_url", "created_at", "updated_at").
		From("labs").
		OrderBy("course_id", "number")

	if courseID != nil {
		builder = builder.Where(sq.Eq{"course_id": *courseID})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get labs: %w", err)
	}
	defer rows.Close()

	var labs []models.Lab
	for rows.Next() {
		var l models.Lab
		if err := rows.Scan(&l.ID, &l.CourseID, &l.Number, &l.Title, &l.Description, &l.Deadline, &l.MaxScore, &l.GithubURL, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error for lab: %w", err)
		}
		labs = append(labs, l)
	}

	return labs, nil
}


func (r *LabRepository) GetByID(id int) (*models.Lab, error) {
	query, args, err := sq.Select("id", "course_id", "number", "title", "description", "deadline", "max_score", "github_url", "created_at", "updated_at").
		From("labs").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var l models.Lab
	err = r.db.QueryRow(query, args...).Scan(&l.ID, &l.CourseID, &l.Number, &l.Title, &l.Description, &l.Deadline, &l.MaxScore, &l.GithubURL, &l.CreatedAt, &l.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get lab: %w", err)
	}

	return &l, nil
}

func (r *LabRepository) Create(courseID, number, maxScore int, title, description, githubURL string, deadline *string) (*models.Lab, error) {
	query, args, _ := sq.Insert("labs").
		Columns("course_id", "number", "title", "description", "deadline", "max_score", "github_url").
		Values(courseID, number, title, description, deadline, maxScore, githubURL).
		ToSql()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	return &models.Lab{
		ID:          int(id),
		CourseID:    courseID,
		Number:      number,
		Title:       title,
		Description: description,
		MaxScore:    maxScore,
		GithubURL:   &githubURL,
	}, nil
}

func (r *LabRepository) Update(id, courseID, number, maxScore int, title, description, githubURL string, deadline *string) error {
	query, args, _ := sq.Update("labs").
		Set("course_id", courseID).
		Set("number", number).
		Set("title", title).
		Set("description", description).
		Set("deadline", deadline).
		Set("max_score", maxScore).
		Set("github_url", githubURL).
		Where(sq.Eq{"id": id}).
		ToSql()

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *LabRepository) Delete(id int) error {
	query, args, _ := sq.Delete("labs").
		Where(sq.Eq{"id": id}).
		ToSql()

	_, err := r.db.Exec(query, args...)
	return err
}
