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
// ticketsPerPage - количество билетов на одну страницу:
//   - 0 = все билеты подряд без разрывов страниц
//   - 1+ = разбивка по страницам с указанным количеством билетов
func (s *DocumentService) GenerateTicketsDocument(tickets []models.Ticket, ticketsPerPage int) []byte {
	var builder strings.Builder

	// Визуальный разделитель страниц (60 символов)
	pageSeparator := strings.Repeat("═", 60)

	for i, ticket := range tickets {
		// Заголовок билета
		builder.WriteString(fmt.Sprintf("Билет № %d\n\n", ticket.Number))

		// Вопросы билета
		for j, question := range ticket.Questions {
			builder.WriteString(fmt.Sprintf("%d. %s\n", j+1, question.Question))
		}

		isLastTicket := i == len(tickets)-1

		if !isLastTicket {
			if ticketsPerPage <= 0 {
				// Режим "все подряд" — только пустая строка между билетами
				builder.WriteString("\n")
			} else {
				// Режим с разбивкой по страницам
				ticketNumberOnPage := (i + 1) % ticketsPerPage

				if ticketNumberOnPage == 0 {
					// Достигли конца страницы — добавляем разделитель и Form Feed
					builder.WriteString("\n")
					builder.WriteString(pageSeparator)
					builder.WriteString("\n")
					builder.WriteString("\f") // Form Feed — разрыв страницы при печати
				} else {
					// Просто отступ между билетами на одной странице
					builder.WriteString("\n")
					builder.WriteString(strings.Repeat("-", 40))
					builder.WriteString("\n\n")
				}
			}
		}
	}

	return []byte(builder.String())
}
