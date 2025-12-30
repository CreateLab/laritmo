package models

import "time"

type GradeSheet struct {
	ID          int       `json:"id" db:"id"`
	CourseID    int       `json:"course_id" db:"course_id"`
	SheetURL    string    `json:"sheet_url" db:"sheet_url"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
