package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/CreateLab/laritmo/internal/repository"
	"github.com/CreateLab/laritmo/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
	logger      *slog.Logger
}

func NewUserHandler(userService *services.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetAll godoc
// @Summary      Get all admin users
// @Description  Get list of all users with admin or owner role
// @Tags         admin-users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/admin/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAllAdmins()
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get users", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	if users == nil {
		users = []*models.User{}
	}

	c.JSON(http.StatusOK, users)
}

// GetByID godoc
// @Summary      Get user by ID
// @Description  Get user details by ID
// @Tags         admin-users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/admin/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "Failed to get user", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Create godoc
// @Summary      Create user
// @Description  Create a new user (admin or student)
// @Tags         admin-users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user  body      models.CreateUserRequest  true  "User data"
// @Success      201   {object}  models.User
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/admin/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "Failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "User created", "id", user.ID, "username", user.Username)
	c.JSON(http.StatusCreated, user)
}

// Update godoc
// @Summary      Update user
// @Description  Update user by ID
// @Tags         admin-users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int                       true  "User ID"
// @Param        user  body      models.UpdateUserRequest  true  "User data"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/admin/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err = h.userService.UpdateUser(id, &req)
	if err != nil {
		if errors.Is(err, services.ErrCannotModifyOwner) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify owner account"})
			return
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "Failed to update user", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "User updated", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// Delete godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         admin-users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/admin/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, services.ErrCannotDeleteOwner) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete owner account"})
			return
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "Failed to delete user", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "User deleted", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// Deactivate godoc
// @Summary      Deactivate user
// @Description  Deactivate user account by ID
// @Tags         admin-users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/admin/users/{id}/deactivate [put]
func (h *UserHandler) Deactivate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.userService.DeactivateUser(id)
	if err != nil {
		if errors.Is(err, services.ErrCannotDeactivateOwner) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot deactivate owner account"})
			return
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "Failed to deactivate user", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "User deactivated", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "User deactivated"})
}

// Activate godoc
// @Summary      Activate user
// @Description  Activate user account by ID
// @Tags         admin-users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/admin/users/{id}/activate [put]
func (h *UserHandler) Activate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.userService.ActivateUser(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "Failed to activate user", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "User activated", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "User activated"})
}

// ResetPassword godoc
// @Summary      Reset user password
// @Description  Reset password for a user (admin action)
// @Tags         admin-users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id        path      int                         true  "User ID"
// @Param        password  body      models.ResetPasswordRequest true  "New password"
// @Success      200       {object}  map[string]string
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      403       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/admin/users/{id}/password [put]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.userService.ResetPassword(adminID.(int), targetUserID, req.NewPassword)
	if err != nil {
		if errors.Is(err, services.ErrCannotResetOwnPass) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Use password change endpoint to update your own password"})
			return
		}
		if errors.Is(err, services.ErrCannotResetOwnerPass) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot reset password for owner account"})
			return
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "Failed to reset password", "error", err, "target_user_id", targetUserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Password reset", "target_user_id", targetUserID, "admin_id", adminID)
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ChangePassword godoc
// @Summary      Change own password
// @Description  Change password for the authenticated user
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        password  body      models.ChangePasswordRequest  true  "Password change data"
// @Success      200       {object}  map[string]string
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/auth/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.userService.ChangePassword(userID.(int), req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, services.ErrInvalidOldPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid old password"})
			return
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "Failed to change password", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Password changed", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
