package models

import "time"

type Course struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Semester    string    `json:"semester" db:"semester"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
