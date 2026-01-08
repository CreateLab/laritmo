package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTicketService - мок для TicketService
type MockTicketService struct {
	mock.Mock
}

func (m *MockTicketService) GenerateRandomTicket(ctx context.Context, courseID int, questionsCount int) (*models.Ticket, error) {
	args := m.Called(ctx, courseID, questionsCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Ticket), args.Error(1)
}

func (m *MockTicketService) GenerateMultipleTickets(ctx context.Context, courseID int, ticketCount, questionsPerTicket int) ([]models.Ticket, error) {
	args := m.Called(ctx, courseID, ticketCount, questionsPerTicket)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Ticket), args.Error(1)
}

// MockDocumentService - мок для DocumentService
type MockDocumentService struct {
	mock.Mock
}

func (m *MockDocumentService) GenerateTicketsDocument(tickets []models.Ticket) []byte {
	args := m.Called(tickets)
	return args.Get(0).([]byte)
}

// MockCourseRepository - мок для CourseRepository
type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) GetByID(id int) (*models.Course, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Course), args.Error(1)
}

func TestTicketHandler_GetRandomTicket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := slog.Default()

	tests := []struct {
		name             string
		courseID         string
		questions        string
		mockCourse       *models.Course
		mockCourseErr    error
		mockTicket       *models.Ticket
		mockTicketErr    error
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "successful generation",
			courseID:  "1",
			questions: "5",
			mockCourse: &models.Course{
				ID:   1,
				Name: "Test Course",
			},
			mockTicket: &models.Ticket{
				Number: 1,
				Questions: []models.Question{
					{Number: 1, Section: "Section A", Question: "Question 1"},
					{Number: 2, Section: "Section B", Question: "Question 2"},
				},
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]models.Ticket
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "ticket")
				assert.Equal(t, 1, response["ticket"].Number)
				assert.Len(t, response["ticket"].Questions, 2)
			},
		},
		{
			name:           "invalid course ID",
			courseID:       "invalid",
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid course ID", response["error"])
			},
		},
		{
			name:           "course not found",
			courseID:       "1",
			questions:      "5",
			mockCourse:     nil,
			expectedStatus: http.StatusNotFound,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Course not found", response["error"])
			},
		},
		{
			name:           "database error getting course",
			courseID:       "1",
			questions:      "5",
			mockCourseErr:  errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to get course", response["error"])
			},
		},
		{
			name:      "invalid questions parameter",
			courseID:  "1",
			questions: "invalid",
			mockCourse: &models.Course{
				ID:   1,
				Name: "Test Course",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid questions parameter", response["error"])
			},
		},
		{
			name:      "questions count too low",
			courseID:  "1",
			questions: "0",
			mockCourse: &models.Course{
				ID:   1,
				Name: "Test Course",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "Questions count must be between")
			},
		},
		{
			name:      "questions count too high",
			courseID:  "1",
			questions: "51",
			mockCourse: &models.Course{
				ID:   1,
				Name: "Test Course",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "Questions count must be between")
			},
		},
		{
			name:      "default questions count",
			courseID:  "1",
			questions: "",
			mockCourse: &models.Course{
				ID:   1,
				Name: "Test Course",
			},
			mockTicket: &models.Ticket{
				Number:    1,
				Questions: []models.Question{},
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]models.Ticket
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "ticket")
			},
		},
		{
			name:      "ticket generation error",
			courseID:  "1",
			questions: "5",
			mockCourse: &models.Course{
				ID:   1,
				Name: "Test Course",
			},
			mockTicketErr:  errors.New("not enough questions"),
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to generate ticket", response["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTicketService := new(MockTicketService)
			mockDocumentService := new(MockDocumentService)
			mockCourseRepo := new(MockCourseRepository)

			// Настройка моков
			if tt.courseID != "invalid" {
				courseIDInt, _ := strconv.Atoi(tt.courseID)
				mockCourseRepo.On("GetByID", courseIDInt).Return(tt.mockCourse, tt.mockCourseErr)

				// Настраиваем мок для генерации билета только если курс существует и нет ошибок валидации
				if tt.mockCourse != nil && tt.mockTicketErr == nil {
					questionsCount := 10 // default
					if tt.questions != "" {
						questionsCount, _ = strconv.Atoi(tt.questions)
					}
					// Не настраиваем мок если questionsCount невалидный (валидация происходит до вызова сервиса)
					if questionsCount >= 1 && questionsCount <= 50 {
						mockTicketService.On("GenerateRandomTicket", mock.Anything, courseIDInt, questionsCount).Return(tt.mockTicket, tt.mockTicketErr)
					}
				} else if tt.mockCourse != nil && tt.mockTicketErr != nil {
					// Если есть ошибка генерации билета, настраиваем мок
					questionsCount := 10 // default
					if tt.questions != "" {
						questionsCount, _ = strconv.Atoi(tt.questions)
					}
					if questionsCount >= 1 && questionsCount <= 50 {
						mockTicketService.On("GenerateRandomTicket", mock.Anything, courseIDInt, questionsCount).Return(tt.mockTicket, tt.mockTicketErr)
					}
				}
			}

			handler := NewTicketHandler(mockTicketService, mockDocumentService, mockCourseRepo, logger)

			router := gin.New()
			router.GET("/courses/:id/tickets/random", handler.GetRandomTicket)

			url := "/courses/" + tt.courseID + "/tickets/random"
			if tt.questions != "" {
				url += "?questions=" + tt.questions
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}

			mockTicketService.AssertExpectations(t)
			mockCourseRepo.AssertExpectations(t)
		})
	}
}

func TestTicketHandler_GenerateTicketsDocument(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := slog.Default()

	tests := []struct {
		name             string
		courseID         string
		requestBody      models.TicketGenerationRequest
		mockCourse       *models.Course
		mockCourseErr    error
		mockTickets      []models.Ticket
		mockTicketsErr   error
		mockDocument     []byte
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:     "successful generation",
			courseID: "1",
			requestBody: models.TicketGenerationRequest{
				QuestionsPerTicket: 5,
				TicketCount:        10,
			},
			mockCourse: &models.Course{
				ID:   1,
				Name: "Test Course",
			},
			mockTickets: []models.Ticket{
				{Number: 1, Questions: []models.Question{{Number: 1, Section: "A", Question: "Q1"}}},
				{Number: 2, Questions: []models.Question{{Number: 2, Section: "B", Question: "Q2"}}},
			},
			mockDocument:   []byte("Билет № 1\n\n1. Q1\n   (раздел: A)\n\n"),
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
				assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
				assert.Contains(t, w.Header().Get("Content-Disposition"), "tickets_test_course.txt")
				assert.Equal(t, "Билет № 1\n\n1. Q1\n   (раздел: A)\n\n", w.Body.String())
			},
		},
		{
			name:           "invalid course ID",
			courseID:       "invalid",
			requestBody:    models.TicketGenerationRequest{QuestionsPerTicket: 5, TicketCount: 10},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid course ID", response["error"])
			},
		},
		{
			name:           "course not found",
			courseID:       "1",
			requestBody:    models.TicketGenerationRequest{QuestionsPerTicket: 5, TicketCount: 10},
			mockCourse:     nil,
			expectedStatus: http.StatusNotFound,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Course not found", response["error"])
			},
		},
		{
			name:           "invalid request body",
			courseID:       "1",
			requestBody:    models.TicketGenerationRequest{}, // пустой body
			mockCourse:     &models.Course{ID: 1, Name: "Test"},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid request format", response["error"])
			},
		},
		{
			name:           "tickets generation error",
			courseID:       "1",
			requestBody:    models.TicketGenerationRequest{QuestionsPerTicket: 5, TicketCount: 10},
			mockCourse:     &models.Course{ID: 1, Name: "Test"},
			mockTicketsErr: errors.New("not enough questions"),
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to generate tickets", response["error"])
			},
		},
		{
			name:     "filename with special characters",
			courseID: "1",
			requestBody: models.TicketGenerationRequest{
				QuestionsPerTicket: 5,
				TicketCount:        10,
			},
			mockCourse: &models.Course{
				ID:   1,
				Name: "Test/Course Name",
			},
			mockTickets: []models.Ticket{
				{Number: 1, Questions: []models.Question{}},
			},
			mockDocument:   []byte("test"),
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				contentDisposition := w.Header().Get("Content-Disposition")
				assert.Contains(t, contentDisposition, "tickets_test_course_name.txt")
				assert.NotContains(t, contentDisposition, "/")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTicketService := new(MockTicketService)
			mockDocumentService := new(MockDocumentService)
			mockCourseRepo := new(MockCourseRepository)

			// Настройка моков
			if tt.courseID != "invalid" {
				courseIDInt, _ := strconv.Atoi(tt.courseID)
				mockCourseRepo.On("GetByID", courseIDInt).Return(tt.mockCourse, tt.mockCourseErr)

				// Настраиваем мок для генерации билетов только если курс существует и body валидный
				if tt.mockCourse != nil {
					// Проверяем валидность body (если requestBody пустой, это ошибка валидации)
					if tt.requestBody.QuestionsPerTicket > 0 && tt.requestBody.TicketCount > 0 {
						mockTicketService.On("GenerateMultipleTickets", mock.Anything, courseIDInt, tt.requestBody.TicketCount, tt.requestBody.QuestionsPerTicket).Return(tt.mockTickets, tt.mockTicketsErr)
						if len(tt.mockTickets) > 0 && tt.mockTicketsErr == nil {
							mockDocumentService.On("GenerateTicketsDocument", tt.mockTickets).Return(tt.mockDocument)
						}
					}
				}
			}

			handler := NewTicketHandler(mockTicketService, mockDocumentService, mockCourseRepo, logger)

			router := gin.New()
			router.POST("/courses/:id/tickets/generate", handler.GenerateTicketsDocument)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/courses/"+tt.courseID+"/tickets/generate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}

			mockTicketService.AssertExpectations(t)
			mockCourseRepo.AssertExpectations(t)
			if len(tt.mockTickets) > 0 && tt.mockTicketsErr == nil {
				mockDocumentService.AssertExpectations(t)
			}
		})
	}
}
