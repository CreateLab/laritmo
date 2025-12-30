package models

import "time"

type ExamQuestion struct {
	ID        int       `json:"id" db:"id"`
	CourseID  int       `json:"course_id" db:"course_id"`
	Number    int       `json:"number" db:"number"`
	Section   string    `json:"section" db:"section"`
	Question  string    `json:"question" db:"question"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
