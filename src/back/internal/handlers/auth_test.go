package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CreateLab/laritmo/internal/auth"
	"github.com/CreateLab/laritmo/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secretKey := "test-secret-key"
	jwtManager := auth.NewJWTManager(secretKey, 24)
	logger := slog.Default()

	t.Run("successful login", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		password := "testpassword123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		user := &models.User{
			ID:           1,
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
			Role:         "admin",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockRepo.On("GetByUsername", "testuser").Return(user, nil)

		handler := NewAuthHandlerWithRepo(mockRepo, jwtManager, logger)

		router := gin.New()
		router.POST("/login", handler.Login)

		reqBody := LoginRequest{
			Username: "testuser",
			Password: password,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response LoginResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, user.ID, response.User.ID)
		assert.Equal(t, user.Username, response.User.Username)
		assert.Equal(t, user.Email, response.User.Email)
		assert.Equal(t, user.Role, response.User.Role)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid username", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByUsername", "nonexistent").Return(nil, errors.New("user not found"))

		handler := NewAuthHandlerWithRepo(mockRepo, jwtManager, logger)

		router := gin.New()
		router.POST("/login", handler.Login)

		reqBody := LoginRequest{
			Username: "nonexistent",
			Password: "password",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid username or password", response["error"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		password := "correctpassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		user := &models.User{
			ID:           1,
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
			Role:         "admin",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockRepo.On("GetByUsername", "testuser").Return(user, nil)

		handler := NewAuthHandlerWithRepo(mockRepo, jwtManager, logger)

		router := gin.New()
		router.POST("/login", handler.Login)

		reqBody := LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid username or password", response["error"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid request format", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		handler := NewAuthHandlerWithRepo(mockRepo, jwtManager, logger)

		router := gin.New()
		router.POST("/login", handler.Login)

		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request format", response["error"])

		mockRepo.AssertNotCalled(t, "GetByUsername")
	})

	t.Run("missing required fields", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		handler := NewAuthHandlerWithRepo(mockRepo, jwtManager, logger)

		router := gin.New()
		router.POST("/login", handler.Login)

		reqBody := LoginRequest{
			Username: "testuser",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request format", response["error"])

		mockRepo.AssertNotCalled(t, "GetByUsername")
	})
}
