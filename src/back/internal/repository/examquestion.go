package repository

import (
	"database/sql"
	"fmt"

	"github.com/CreateLab/laritmo/internal/models"
	sq "github.com/Masterminds/squirrel"
)

type ExamQuestionRepository struct {
	db *sql.DB
}

func NewExamQuestionRepository(db *sql.DB) *ExamQuestionRepository {
	return &ExamQuestionRepository{db: db}
}


func (r *ExamQuestionRepository) GetAll() ([]models.ExamQuestion, error) {
	query, args, err := sq.Select("id", "course_id", "number", "section", "question", "created_at", "updated_at").
		From("exam_questions").
		OrderBy("course_id", "section", "number").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam questions: %w", err)
	}
	defer rows.Close()

	var questions []models.ExamQuestion
	for rows.Next() {
		var q models.ExamQuestion
		if err := rows.Scan(&q.ID, &q.CourseID, &q.Number, &q.Section, &q.Question, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error exam question: %w", err)
		}
		questions = append(questions, q)
	}

	if questions == nil {
		questions = []models.ExamQuestion{}
	}

	return questions, nil
}


func (r *ExamQuestionRepository) GetByCourseID(courseID int) ([]models.ExamQuestion, error) {
	query, args, err := sq.Select("id", "course_id", "number", "section", "question", "created_at", "updated_at").
		From("exam_questions").
		Where(sq.Eq{"course_id": courseID}).
		OrderBy("section ASC", "number ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam questions: %w", err)
	}
	defer rows.Close()

	var questions []models.ExamQuestion
	for rows.Next() {
		var q models.ExamQuestion
		if err := rows.Scan(&q.ID, &q.CourseID, &q.Number, &q.Section, &q.Question, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error exam question: %w", err)
		}
		questions = append(questions, q)
	}

	if questions == nil {
		questions = []models.ExamQuestion{}
	}

	return questions, nil
}


func (r *ExamQuestionRepository) GetByID(id int) (*models.ExamQuestion, error) {
	query, args, err := sq.Select("id", "course_id", "number", "section", "question", "created_at", "updated_at").
		From("exam_questions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var q models.ExamQuestion
	err = r.db.QueryRow(query, args...).Scan(&q.ID, &q.CourseID, &q.Number, &q.Section, &q.Question, &q.CreatedAt, &q.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get exam question: %w", err)
	}

	return &q, nil
}


func (r *ExamQuestionRepository) Create(courseID, number int, section, question string) (*models.ExamQuestion, error) {
	query, args, err := sq.Insert("exam_questions").
		Columns("course_id", "number", "section", "question").
		Values(courseID, number, section, question).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create exam question: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get ID: %w", err)
	}

	
	examQuestion, err := r.GetByID(int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get created exam question: %w", err)
	}
	if examQuestion == nil {
		return nil, fmt.Errorf("created exam question not found")
	}

	return examQuestion, nil
}


func (r *ExamQuestionRepository) Update(id, courseID, number int, section, question string) error {
	query, args, err := sq.Update("exam_questions").
		Set("course_id", courseID).
		Set("number", number).
		Set("section", section).
		Set("question", question).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update exam question: %w", err)
	}

	return nil
}


func (r *ExamQuestionRepository) Delete(id int) error {
	query, args, err := sq.Delete("exam_questions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete exam question: %w", err)
	}

	return nil
}


func (r *ExamQuestionRepository) BulkCreate(questions []models.ExamQuestion) error {
	if len(questions) == 0 {
		return nil
	}

	builder := sq.Insert("exam_questions").
		Columns("course_id", "number", "section", "question")

	for _, q := range questions {
		builder = builder.Values(q.CourseID, q.Number, q.Section, q.Question)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk create exam questions: %w", err)
	}

	return nil
}


func (r *ExamQuestionRepository) DeleteByCourseID(courseID int) error {
	query, args, err := sq.Delete("exam_questions").
		Where(sq.Eq{"course_id": courseID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete exam questions by course_id: %w", err)
	}

	return nil
}
