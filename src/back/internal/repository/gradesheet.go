package repository

import (
	"database/sql"
	"fmt"

	"github.com/CreateLab/laritmo/internal/models"
	sq "github.com/Masterminds/squirrel"
)

type GradeSheetRepository struct {
	db *sql.DB
}

func NewGradeSheetRepository(db *sql.DB) *GradeSheetRepository {
	return &GradeSheetRepository{db: db}
}


func (r *GradeSheetRepository) GetAll(courseID *int) ([]models.GradeSheet, error) {
	builder := sq.Select("id", "course_id", "sheet_url", "description", "created_at", "updated_at").
		From("grade_sheets").
		OrderBy("course_id")

	if courseID != nil {
		builder = builder.Where(sq.Eq{"course_id": *courseID})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get grade sheets: %w", err)
	}
	defer rows.Close()

	var sheets []models.GradeSheet
	for rows.Next() {
		var s models.GradeSheet
		if err := rows.Scan(&s.ID, &s.CourseID, &s.SheetURL, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error grade sheet: %w", err)
		}
		sheets = append(sheets, s)
	}

	return sheets, nil
}


func (r *GradeSheetRepository) GetByID(id int) (*models.GradeSheet, error) {
	query, args, err := sq.Select("id", "course_id", "sheet_url", "description", "created_at", "updated_at").
		From("grade_sheets").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var s models.GradeSheet
	err = r.db.QueryRow(query, args...).Scan(&s.ID, &s.CourseID, &s.SheetURL, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get grade sheet: %w", err)
	}

	return &s, nil
}


func (r *GradeSheetRepository) Create(courseID int, sheetURL, description string) (*models.GradeSheet, error) {
	query, args, err := sq.Insert("grade_sheets").
		Columns("course_id", "sheet_url", "description").
		Values(courseID, sheetURL, description).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grade sheet: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get ID: %w", err)
	}

	
	gradeSheet, err := r.GetByID(int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get created grade sheet: %w", err)
	}
	if gradeSheet == nil {
		return nil, fmt.Errorf("created grade sheet not found")
	}

	return gradeSheet, nil
}


func (r *GradeSheetRepository) Update(id int, sheetURL, description string) error {
	query, args, err := sq.Update("grade_sheets").
		Set("sheet_url", sheetURL).
		Set("description", description).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update grade sheet: %w", err)
	}

	return nil
}


func (r *GradeSheetRepository) Delete(id int) error {
	query, args, err := sq.Delete("grade_sheets").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete grade sheet: %w", err)
	}

	return nil
}
