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

// GenerateTicketsDocument generates TXT document with tickets in memory
// ticketsPerPage - number of tickets per page:
//   - 0 = all tickets in sequence without page breaks
//   - 1+ = pagination with specified number of tickets
func (s *DocumentService) GenerateTicketsDocument(tickets []models.Ticket, ticketsPerPage int) []byte {
	var builder strings.Builder

	// Visual page separator (60 characters)
	pageSeparator := strings.Repeat("═", 60)

	for i, ticket := range tickets {
		// Ticket header
		builder.WriteString(fmt.Sprintf("Ticket #%d\n\n", ticket.Number))

		// Ticket questions
		for j, question := range ticket.Questions {
			builder.WriteString(fmt.Sprintf("%d. %s\n", j+1, question.Question))
		}

		isLastTicket := i == len(tickets)-1

		if !isLastTicket {
			if ticketsPerPage <= 0 {
				// Sequential mode - only blank line between tickets
				builder.WriteString("\n")
			} else {
				// Pagination mode
				ticketNumberOnPage := (i + 1) % ticketsPerPage

				if ticketNumberOnPage == 0 {
					// Reached end of page - add separator and Form Feed
					builder.WriteString("\n")
					builder.WriteString(pageSeparator)
					builder.WriteString("\n")
					builder.WriteString("\f") // Form Feed - page break for printing
				} else {
					// Just spacing between tickets on same page
					builder.WriteString("\n")
					builder.WriteString(strings.Repeat("-", 40))
					builder.WriteString("\n\n")
				}
			}
		}
	}

	return []byte(builder.String())
}
