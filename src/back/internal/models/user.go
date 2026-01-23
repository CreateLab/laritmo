package models

import "time"

type User struct {
	ID           int        `json:"id"`
	Email        string     `json:"email"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"` // Excluded from JSON response
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// User roles
const (
	RoleOwner   = "owner"
	RoleAdmin   = "admin"
	RoleStudent = "student"
)
