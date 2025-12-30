package models

import "time"

type Lab struct {
	ID          int        `json:"id" db:"id"`
	CourseID    int        `json:"course_id" db:"course_id"`
	Number      int        `json:"number" db:"number"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Deadline    *time.Time `json:"deadline,omitempty" db:"deadline"`
	MaxScore    int        `json:"max_score" db:"max_score"`
	GithubURL   *string    `json:"github_url,omitempty" db:"github_url"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
