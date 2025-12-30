package handlers

import (
	"log/slog"
	"net/http"

	"github.com/CreateLab/laritmo/internal/auth"
	"github.com/CreateLab/laritmo/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo   *repository.UserRepository
	jwtManager *auth.JWTManager
	logger     *slog.Logger
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtManager *auth.JWTManager, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user"`
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest   true  "Login credentials"
// @Success      200          {object}  LoginResponse
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "User not found", "username", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Invalid password", "username", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := h.jwtManager.GenerateToken(user)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Token generation error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization error"})
		return
	}

	response := LoginResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.Email = user.Email
	response.User.Role = user.Role

	h.logger.InfoContext(c.Request.Context(), "Login successful", "username", user.Username)
	c.JSON(http.StatusOK, response)
}
