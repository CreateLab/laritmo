package models

import "time"

type Lecture struct {
	ID        int       `json:"id" db:"id"`
	CourseID  int       `json:"course_id" db:"course_id"`
	Week      int       `json:"week" db:"week"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	GithubURL *string   `json:"github_url,omitempty" db:"github_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
