package services

import (
	"strings"
	"testing"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestDocumentService_GenerateTicketsDocument(t *testing.T) {
	service := NewDocumentService()

	tests := []struct {
		name           string
		tickets        []models.Ticket
		validateOutput func(*testing.T, []byte)
	}{
		{
			name: "single ticket",
			tickets: []models.Ticket{
				{
					Number: 1,
					Questions: []models.Question{
						{Number: 1, Section: "Основы", Question: "Что такое Go?"},
						{Number: 2, Section: "Продвинутое", Question: "Что такое горутины?"},
					},
				},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				assert.Contains(t, text, "Билет № 1")
				assert.Contains(t, text, "1. Что такое Go?")
				assert.Contains(t, text, "(раздел: Основы)")
				assert.Contains(t, text, "2. Что такое горутины?")
				assert.Contains(t, text, "(раздел: Продвинутое)")
			},
		},
		{
			name: "multiple tickets",
			tickets: []models.Ticket{
				{
					Number: 1,
					Questions: []models.Question{
						{Number: 1, Section: "Раздел A", Question: "Вопрос 1"},
					},
				},
				{
					Number: 2,
					Questions: []models.Question{
						{Number: 2, Section: "Раздел B", Question: "Вопрос 2"},
					},
				},
				{
					Number: 3,
					Questions: []models.Question{
						{Number: 3, Section: "Раздел C", Question: "Вопрос 3"},
					},
				},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				assert.Contains(t, text, "Билет № 1")
				assert.Contains(t, text, "Билет № 2")
				assert.Contains(t, text, "Билет № 3")
				assert.Contains(t, text, "Вопрос 1")
				assert.Contains(t, text, "Вопрос 2")
				assert.Contains(t, text, "Вопрос 3")
			},
		},
		{
			name:    "empty tickets list",
			tickets: []models.Ticket{},
			validateOutput: func(t *testing.T, output []byte) {
				assert.Empty(t, output)
			},
		},
		{
			name: "ticket with single question",
			tickets: []models.Ticket{
				{
					Number: 1,
					Questions: []models.Question{
						{Number: 1, Section: "Основы", Question: "Единственный вопрос"},
					},
				},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				assert.Contains(t, text, "Билет № 1")
				assert.Contains(t, text, "1. Единственный вопрос")
				assert.Contains(t, text, "(раздел: Основы)")
			},
		},
		{
			name: "UTF-8 encoding with Russian characters",
			tickets: []models.Ticket{
				{
					Number: 1,
					Questions: []models.Question{
						{Number: 1, Section: "Основы программирования", Question: "Что такое переменная в программировании?"},
						{Number: 2, Section: "Алгоритмы и структуры данных", Question: "Объясните принцип работы стека."},
					},
				},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				// Проверяем, что русские символы корректно отображаются
				assert.Contains(t, text, "Основы программирования")
				assert.Contains(t, text, "Алгоритмы и структуры данных")
				assert.Contains(t, text, "Что такое переменная")
				assert.Contains(t, text, "принцип работы стека")

				// Проверяем структуру
				lines := strings.Split(text, "\n")
				assert.Contains(t, lines[0], "Билет № 1")
				assert.Contains(t, lines[2], "1. Что такое переменная")
				assert.Contains(t, lines[3], "раздел: Основы программирования")
			},
		},
		{
			name: "correct numbering of questions",
			tickets: []models.Ticket{
				{
					Number: 1,
					Questions: []models.Question{
						{Number: 5, Section: "A", Question: "Question 5"},
						{Number: 10, Section: "B", Question: "Question 10"},
						{Number: 15, Section: "C", Question: "Question 15"},
					},
				},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				// В билете вопросы должны быть пронумерованы 1, 2, 3 (не 5, 10, 15)
				assert.Contains(t, text, "1. Question 5")
				assert.Contains(t, text, "2. Question 10")
				assert.Contains(t, text, "3. Question 15")
				// Но оригинальный номер вопроса не должен быть в тексте билета
				assert.NotContains(t, text, "5. Question 5")
			},
		},
		{
			name: "correct ticket numbering",
			tickets: []models.Ticket{
				{Number: 1, Questions: []models.Question{{Number: 1, Section: "A", Question: "Q1"}}},
				{Number: 2, Questions: []models.Question{{Number: 2, Section: "B", Question: "Q2"}}},
				{Number: 3, Questions: []models.Question{{Number: 3, Section: "C", Question: "Q3"}}},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				// Проверяем, что номера билетов правильные
				assert.Contains(t, text, "Билет № 1")
				assert.Contains(t, text, "Билет № 2")
				assert.Contains(t, text, "Билет № 3")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := service.GenerateTicketsDocument(tt.tickets)
			tt.validateOutput(t, output)
		})
	}
}

func TestDocumentService_FormatStructure(t *testing.T) {
	service := NewDocumentService()

	ticket := models.Ticket{
		Number: 1,
		Questions: []models.Question{
			{Number: 1, Section: "Section A", Question: "Question 1"},
			{Number: 2, Section: "Section B", Question: "Question 2"},
		},
	}

	output := service.GenerateTicketsDocument([]models.Ticket{ticket})
	text := string(output)
	lines := strings.Split(text, "\n")

	// Проверяем структуру:
	// Билет № 1
	//
	// 1. Question 1
	//    (раздел: Section A)
	//
	// 2. Question 2
	//    (раздел: Section B)

	assert.Contains(t, lines[0], "Билет № 1")
	assert.True(t, lines[1] == "", "line 1 should be empty")
	assert.Contains(t, lines[2], "1. Question 1")
	assert.Contains(t, lines[3], "(раздел: Section A)")
	assert.Contains(t, lines[5], "2. Question 2")
	assert.Contains(t, lines[6], "(раздел: Section B)")
}
