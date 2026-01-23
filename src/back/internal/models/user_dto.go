package models

// CreateUserRequest - DTO for creating a new user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=admin student"`
}

// UpdateUserRequest - DTO for updating a user
type UpdateUserRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Role     *string `json:"role" binding:"omitempty,oneof=owner admin student"`
	IsActive *bool   `json:"is_active"`
}

// ChangePasswordRequest - DTO for user changing their own password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPasswordRequest - DTO for admin resetting user password
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
