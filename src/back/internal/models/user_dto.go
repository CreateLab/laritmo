package models

// CreateUserRequest - DTO для создания нового пользователя
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=admin student"`
}

// UpdateUserRequest - DTO для обновления пользователя
type UpdateUserRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Role     *string `json:"role" binding:"omitempty,oneof=owner admin student"`
	IsActive *bool   `json:"is_active"`
}

// ChangePasswordRequest - DTO для смены пароля самим пользователем
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPasswordRequest - DTO для сброса пароля администратором
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
