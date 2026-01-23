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
						{Number: 1, Section: "Basics", Question: "What is Go?"},
						{Number: 2, Section: "Advanced", Question: "What are goroutines?"},
					},
				},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				assert.Contains(t, text, "Ticket #1")
				assert.Contains(t, text, "1. What is Go?")
				assert.Contains(t, text, "2. What are goroutines?")
			},
		},
		{
			name: "multiple tickets",
			tickets: []models.Ticket{
				{
					Number: 1,
					Questions: []models.Question{
						{Number: 1, Section: "Section A", Question: "Question 1"},
					},
				},
				{
					Number: 2,
					Questions: []models.Question{
						{Number: 2, Section: "Section B", Question: "Question 2"},
					},
				},
				{
					Number: 3,
					Questions: []models.Question{
						{Number: 3, Section: "Section C", Question: "Question 3"},
					},
				},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				assert.Contains(t, text, "Ticket #1")
				assert.Contains(t, text, "Ticket #2")
				assert.Contains(t, text, "Ticket #3")
				assert.Contains(t, text, "Question 1")
				assert.Contains(t, text, "Question 2")
				assert.Contains(t, text, "Question 3")
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
						{Number: 1, Section: "Basics", Question: "Single question"},
					},
				},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				assert.Contains(t, text, "Ticket #1")
				assert.Contains(t, text, "1. Single question")
			},
		},
		{
			name: "UTF-8 encoding with special characters",
			tickets: []models.Ticket{
				{
					Number: 1,
					Questions: []models.Question{
						{Number: 1, Section: "Programming basics", Question: "What is a variable in programming?"},
						{Number: 2, Section: "Algorithms and data structures", Question: "Explain how a stack works."},
					},
				},
			},
			validateOutput: func(t *testing.T, output []byte) {
				text := string(output)
				// Check that characters are correctly displayed
				assert.Contains(t, text, "What is a variable")
				assert.Contains(t, text, "how a stack works")

				// Check structure
				lines := strings.Split(text, "\n")
				assert.Contains(t, lines[0], "Ticket #1")
				assert.Contains(t, lines[2], "1. What is a variable")
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
				// In ticket, questions should be numbered 1, 2, 3 (not 5, 10, 15)
				assert.Contains(t, text, "1. Question 5")
				assert.Contains(t, text, "2. Question 10")
				assert.Contains(t, text, "3. Question 15")
				// Original question number should not be in ticket text
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
				// Check that ticket numbers are correct
				assert.Contains(t, text, "Ticket #1")
				assert.Contains(t, text, "Ticket #2")
				assert.Contains(t, text, "Ticket #3")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := service.GenerateTicketsDocument(tt.tickets, 1)
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

	output := service.GenerateTicketsDocument([]models.Ticket{ticket}, 1)
	text := string(output)
	lines := strings.Split(text, "\n")

	// Check structure:
	// Ticket #1
	//
	// 1. Question 1
	//
	// 2. Question 2

	assert.Contains(t, lines[0], "Ticket #1")
	assert.True(t, lines[1] == "", "line 1 should be empty")
	assert.Contains(t, lines[2], "1. Question 1")
	assert.Contains(t, lines[3], "2. Question 2")
}

func TestDocumentService_PageBreaks(t *testing.T) {
	service := NewDocumentService()

	tickets := []models.Ticket{
		{Number: 1, Questions: []models.Question{{Number: 1, Section: "A", Question: "Q1"}}},
		{Number: 2, Questions: []models.Question{{Number: 2, Section: "B", Question: "Q2"}}},
		{Number: 3, Questions: []models.Question{{Number: 3, Section: "C", Question: "Q3"}}},
		{Number: 4, Questions: []models.Question{{Number: 4, Section: "D", Question: "Q4"}}},
	}

	t.Run("all in a row (ticketsPerPage=0)", func(t *testing.T) {
		output := service.GenerateTicketsDocument(tickets, 0)
		text := string(output)
		// Should not have page breaks and separators
		formFeedCount := strings.Count(text, "\f")
		assert.Equal(t, 0, formFeedCount)
		assert.NotContains(t, text, "═")
		assert.NotContains(t, text, "----")
		// Should have all tickets
		assert.Contains(t, text, "Ticket #1")
		assert.Contains(t, text, "Ticket #4")
	})

	t.Run("1 ticket per page", func(t *testing.T) {
		output := service.GenerateTicketsDocument(tickets, 1)
		text := string(output)
		// Should have 3 page breaks (after tickets 1, 2, 3)
		formFeedCount := strings.Count(text, "\f")
		assert.Equal(t, 3, formFeedCount)
		// Should have visual separators
		assert.Contains(t, text, "═")
	})

	t.Run("2 tickets per page", func(t *testing.T) {
		output := service.GenerateTicketsDocument(tickets, 2)
		text := string(output)
		// Should have 1 page break (after ticket 2)
		formFeedCount := strings.Count(text, "\f")
		assert.Equal(t, 1, formFeedCount)
		// Should have separators between tickets on same page
		assert.Contains(t, text, "----")
	})

	t.Run("all tickets on one page (ticketsPerPage > count)", func(t *testing.T) {
		output := service.GenerateTicketsDocument(tickets, 10)
		text := string(output)
		// Should not have page breaks
		formFeedCount := strings.Count(text, "\f")
		assert.Equal(t, 0, formFeedCount)
		// But should have separators between tickets
		assert.Contains(t, text, "----")
	})
}
