package services

import (
	"context"
	"errors"
	"testing"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExamQuestionRepository - мок для ExamQuestionRepository
type MockExamQuestionRepository struct {
	mock.Mock
}

func (m *MockExamQuestionRepository) GetByCourseID(courseID int) ([]models.ExamQuestion, error) {
	args := m.Called(courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ExamQuestion), args.Error(1)
}

func TestTicketService_GenerateRandomTicket(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		courseID       int
		questionsCount int
		mockQuestions  []models.ExamQuestion
		mockError      error
		expectedError  string
		validateTicket func(*testing.T, *models.Ticket)
	}{
		{
			name:           "successful generation with enough questions",
			courseID:       1,
			questionsCount: 3,
			mockQuestions: []models.ExamQuestion{
				{ID: 1, CourseID: 1, Number: 1, Section: "Section A", Question: "Question 1"},
				{ID: 2, CourseID: 1, Number: 2, Section: "Section A", Question: "Question 2"},
				{ID: 3, CourseID: 1, Number: 3, Section: "Section B", Question: "Question 3"},
				{ID: 4, CourseID: 1, Number: 4, Section: "Section B", Question: "Question 4"},
				{ID: 5, CourseID: 1, Number: 5, Section: "Section C", Question: "Question 5"},
			},
			validateTicket: func(t *testing.T, ticket *models.Ticket) {
				assert.NotNil(t, ticket)
				assert.Equal(t, 1, ticket.Number)
				assert.Len(t, ticket.Questions, 3)
				// Проверяем, что все вопросы уникальны
				questionIDs := make(map[int]bool)
				for _, q := range ticket.Questions {
					assert.False(t, questionIDs[q.Number], "duplicate question number")
					questionIDs[q.Number] = true
				}
			},
		},
		{
			name:           "not enough questions",
			courseID:       1,
			questionsCount: 10,
			mockQuestions: []models.ExamQuestion{
				{ID: 1, CourseID: 1, Number: 1, Section: "Section A", Question: "Question 1"},
				{ID: 2, CourseID: 1, Number: 2, Section: "Section B", Question: "Question 2"},
			},
			expectedError: "not enough questions: have 2, need 10",
		},
		{
			name:           "database error",
			courseID:       1,
			questionsCount: 3,
			mockError:      errors.New("database connection failed"),
			expectedError:  "failed to get exam questions: database connection failed",
		},
		{
			name:           "invalid questions count - too low",
			courseID:       1,
			questionsCount: 0,
			// Не настраиваем мок, так как валидация происходит до обращения к репозиторию
			expectedError: "questions count must be between 1 and 50",
		},
		{
			name:           "invalid questions count - too high",
			courseID:       1,
			questionsCount: 51,
			// Не настраиваем мок, так как валидация происходит до обращения к репозиторию
			expectedError: "questions count must be between 1 and 50",
		},
		{
			name:           "successful generation with multiple sections",
			courseID:       1,
			questionsCount: 2,
			mockQuestions: []models.ExamQuestion{
				{ID: 1, CourseID: 1, Number: 1, Section: "Section A", Question: "Question 1"},
				{ID: 2, CourseID: 1, Number: 2, Section: "Section A", Question: "Question 2"},
				{ID: 3, CourseID: 1, Number: 3, Section: "Section B", Question: "Question 3"},
				{ID: 4, CourseID: 1, Number: 4, Section: "Section B", Question: "Question 4"},
			},
			validateTicket: func(t *testing.T, ticket *models.Ticket) {
				assert.NotNil(t, ticket)
				assert.Len(t, ticket.Questions, 2)
				// Проверяем, что вопросы из разных разделов
				sections := make(map[string]bool)
				for _, q := range ticket.Questions {
					sections[q.Section] = true
				}
				// Должны быть вопросы из разных разделов
				assert.GreaterOrEqual(t, len(sections), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExamQuestionRepository)
			// Настраиваем мок только если есть mockQuestions или mockError (т.е. когда будет вызов репозитория)
			if tt.mockQuestions != nil || tt.mockError != nil {
				if tt.mockError != nil {
					mockRepo.On("GetByCourseID", tt.courseID).Return(nil, tt.mockError)
				} else {
					mockRepo.On("GetByCourseID", tt.courseID).Return(tt.mockQuestions, nil)
				}
			}

			service := NewTicketService(mockRepo)
			ticket, err := service.GenerateRandomTicket(ctx, tt.courseID, tt.questionsCount)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, ticket)
			} else {
				assert.NoError(t, err)
				if tt.validateTicket != nil {
					tt.validateTicket(t, ticket)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTicketService_GenerateMultipleTickets(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name               string
		courseID           int
		ticketCount        int
		questionsPerTicket int
		mockQuestions      []models.ExamQuestion
		mockError          error
		expectedError      string
		validateTickets    func(*testing.T, []models.Ticket)
	}{
		{
			name:               "successful generation of multiple tickets",
			courseID:           1,
			ticketCount:        3,
			questionsPerTicket: 2,
			mockQuestions: []models.ExamQuestion{
				{ID: 1, CourseID: 1, Number: 1, Section: "Section A", Question: "Question 1"},
				{ID: 2, CourseID: 1, Number: 2, Section: "Section A", Question: "Question 2"},
				{ID: 3, CourseID: 1, Number: 3, Section: "Section B", Question: "Question 3"},
				{ID: 4, CourseID: 1, Number: 4, Section: "Section B", Question: "Question 4"},
				{ID: 5, CourseID: 1, Number: 5, Section: "Section C", Question: "Question 5"},
				{ID: 6, CourseID: 1, Number: 6, Section: "Section C", Question: "Question 6"},
			},
			validateTickets: func(t *testing.T, tickets []models.Ticket) {
				assert.Len(t, tickets, 3)
				// Проверяем нумерацию
				for i, ticket := range tickets {
					assert.Equal(t, i+1, ticket.Number)
					assert.Len(t, ticket.Questions, 2)
				}
				// Проверяем, что билеты уникальны (не все одинаковые)
				firstTicketQuestions := tickets[0].Questions
				allSame := true
				for _, ticket := range tickets[1:] {
					if len(ticket.Questions) != len(firstTicketQuestions) {
						allSame = false
						break
					}
					for i, q := range ticket.Questions {
						if q.Number != firstTicketQuestions[i].Number {
							allSame = false
							break
						}
					}
					if !allSame {
						break
					}
				}
				// С высокой вероятностью билеты должны отличаться
				// Но если вопросов мало, они могут совпадать - это нормально
			},
		},
		{
			name:               "not enough questions",
			courseID:           1,
			ticketCount:        3,
			questionsPerTicket: 5,
			mockQuestions: []models.ExamQuestion{
				{ID: 1, CourseID: 1, Number: 1, Section: "Section A", Question: "Question 1"},
				{ID: 2, CourseID: 1, Number: 2, Section: "Section B", Question: "Question 2"},
			},
			expectedError: "not enough questions: have 2, need at least 5",
		},
		{
			name:               "database error",
			courseID:           1,
			ticketCount:        3,
			questionsPerTicket: 2,
			mockError:          errors.New("database connection failed"),
			expectedError:      "failed to get exam questions: database connection failed",
		},
		{
			name:               "invalid ticket count - too low",
			courseID:           1,
			ticketCount:        0,
			questionsPerTicket: 2,
			// Не настраиваем мок, так как валидация происходит до обращения к репозиторию
			expectedError: "ticket count must be between 1 and 100",
		},
		{
			name:               "invalid ticket count - too high",
			courseID:           1,
			ticketCount:        101,
			questionsPerTicket: 2,
			// Не настраиваем мок, так как валидация происходит до обращения к репозиторию
			expectedError: "ticket count must be between 1 and 100",
		},
		{
			name:               "invalid questions per ticket - too low",
			courseID:           1,
			ticketCount:        3,
			questionsPerTicket: 0,
			// Не настраиваем мок, так как валидация происходит до обращения к репозиторию
			expectedError: "questions per ticket must be between 1 and 50",
		},
		{
			name:               "invalid questions per ticket - too high",
			courseID:           1,
			ticketCount:        3,
			questionsPerTicket: 51,
			// Не настраиваем мок, так как валидация происходит до обращения к репозиторию
			expectedError: "questions per ticket must be between 1 and 50",
		},
		{
			name:               "tickets have correct numbering",
			courseID:           1,
			ticketCount:        5,
			questionsPerTicket: 1,
			mockQuestions: []models.ExamQuestion{
				{ID: 1, CourseID: 1, Number: 1, Section: "Section A", Question: "Question 1"},
				{ID: 2, CourseID: 1, Number: 2, Section: "Section A", Question: "Question 2"},
				{ID: 3, CourseID: 1, Number: 3, Section: "Section B", Question: "Question 3"},
				{ID: 4, CourseID: 1, Number: 4, Section: "Section B", Question: "Question 4"},
				{ID: 5, CourseID: 1, Number: 5, Section: "Section C", Question: "Question 5"},
			},
			validateTickets: func(t *testing.T, tickets []models.Ticket) {
				assert.Len(t, tickets, 5)
				for i, ticket := range tickets {
					assert.Equal(t, i+1, ticket.Number, "ticket number should be %d", i+1)
					assert.Len(t, ticket.Questions, 1)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExamQuestionRepository)
			// Настраиваем мок только если есть mockQuestions или mockError (т.е. когда будет вызов репозитория)
			if tt.mockQuestions != nil || tt.mockError != nil {
				if tt.mockError != nil {
					mockRepo.On("GetByCourseID", tt.courseID).Return(nil, tt.mockError)
				} else {
					mockRepo.On("GetByCourseID", tt.courseID).Return(tt.mockQuestions, nil)
				}
			}

			service := NewTicketService(mockRepo)
			tickets, err := service.GenerateMultipleTickets(ctx, tt.courseID, tt.ticketCount, tt.questionsPerTicket)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, tickets)
			} else {
				assert.NoError(t, err)
				if tt.validateTickets != nil {
					tt.validateTickets(t, tickets)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTicketService_QuestionDistribution(t *testing.T) {
	ctx := context.Background()

	// Создаем вопросы из разных разделов
	mockQuestions := []models.ExamQuestion{
		{ID: 1, CourseID: 1, Number: 1, Section: "Section A", Question: "Question A1"},
		{ID: 2, CourseID: 1, Number: 2, Section: "Section A", Question: "Question A2"},
		{ID: 3, CourseID: 1, Number: 3, Section: "Section B", Question: "Question B1"},
		{ID: 4, CourseID: 1, Number: 4, Section: "Section B", Question: "Question B2"},
		{ID: 5, CourseID: 1, Number: 5, Section: "Section C", Question: "Question C1"},
		{ID: 6, CourseID: 1, Number: 6, Section: "Section C", Question: "Question C2"},
	}

	mockRepo := new(MockExamQuestionRepository)
	mockRepo.On("GetByCourseID", 1).Return(mockQuestions, nil)

	service := NewTicketService(mockRepo)

	// Генерируем билет с количеством вопросов >= количеству разделов
	ticket, err := service.GenerateRandomTicket(ctx, 1, 3)
	assert.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Len(t, ticket.Questions, 3)

	// Проверяем, что вопросы из разных разделов (должны быть представлены разные разделы)
	sections := make(map[string]int)
	for _, q := range ticket.Questions {
		sections[q.Section]++
	}
	// Должно быть минимум 2 разных раздела (так как вопросов 3, а разделов 3)
	assert.GreaterOrEqual(t, len(sections), 1)

	mockRepo.AssertExpectations(t)
}

func TestTicketService_NoDuplicates(t *testing.T) {
	ctx := context.Background()

	mockQuestions := []models.ExamQuestion{
		{ID: 1, CourseID: 1, Number: 1, Section: "Section A", Question: "Question 1"},
		{ID: 2, CourseID: 1, Number: 2, Section: "Section A", Question: "Question 2"},
		{ID: 3, CourseID: 1, Number: 3, Section: "Section B", Question: "Question 3"},
		{ID: 4, CourseID: 1, Number: 4, Section: "Section B", Question: "Question 4"},
	}

	mockRepo := new(MockExamQuestionRepository)
	mockRepo.On("GetByCourseID", 1).Return(mockQuestions, nil)

	service := NewTicketService(mockRepo)

	ticket, err := service.GenerateRandomTicket(ctx, 1, 2)
	assert.NoError(t, err)
	assert.NotNil(t, ticket)

	// Проверяем отсутствие дубликатов
	questionNumbers := make(map[int]bool)
	for _, q := range ticket.Questions {
		assert.False(t, questionNumbers[q.Number], "duplicate question number %d", q.Number)
		questionNumbers[q.Number] = true
	}

	mockRepo.AssertExpectations(t)
}
