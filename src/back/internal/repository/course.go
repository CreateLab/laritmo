package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/CreateLab/laritmo/internal/models"
	sq "github.com/Masterminds/squirrel"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) GetAll() ([]models.Course, error) {
	query, args, err := sq.Select("id", "name", "semester", "description", "created_at", "updated_at").
		From("courses").
		OrderBy("semester", "name").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.ID, &c.Name, &c.Semester, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error for course: %w", err)
		}
		courses = append(courses, c)
	}

	return courses, nil
}

func (r *CourseRepository) GetByID(id int) (*models.Course, error) {
	query, args, err := sq.Select("id", "name", "semester", "description", "created_at", "updated_at").
		From("courses").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var c models.Course
	err = r.db.QueryRow(query, args...).Scan(&c.ID, &c.Name, &c.Semester, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	return &c, nil
}

func (r *CourseRepository) Create(name, semester, description string) (*models.Course, error) {
	query, args, _ := sq.Insert("courses").
		Columns("name", "semester", "description").
		Values(name, semester, description).
		ToSql()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	return &models.Course{
		ID:          int(id),
		Name:        name,
		Semester:    semester,
		Description: description,
	}, nil
}

func (r *CourseRepository) Update(id int, name, semester, description string) error {
	query, args, _ := sq.Update("courses").
		Set("name", name).
		Set("semester", semester).
		Set("description", description).
		Where(sq.Eq{"id": id}).
		ToSql()

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *CourseRepository) Delete(id int) error {
	query, args, _ := sq.Delete("courses").
		Where(sq.Eq{"id": id}).
		ToSql()

	_, err := r.db.Exec(query, args...)
	return err
}
