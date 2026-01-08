package services

import (
	"fmt"
	"strings"

	"github.com/CreateLab/laritmo/internal/models"
)

type DocumentService struct{}

func NewDocumentService() *DocumentService {
	return &DocumentService{}
}

// GenerateTicketsDocument генерирует TXT документ с билетами в памяти
func (s *DocumentService) GenerateTicketsDocument(tickets []models.Ticket) []byte {
	var builder strings.Builder

	for i, ticket := range tickets {
		// Заголовок билета
		builder.WriteString(fmt.Sprintf("Билет № %d\n\n", ticket.Number))

		// Вопросы билета
		for j, question := range ticket.Questions {
			builder.WriteString(fmt.Sprintf("%d. %s\n", j+1, question.Question))
		}

		if i < len(tickets)-1 {
			builder.WriteString("\n")
		}
	}

	return []byte(builder.String())
}
