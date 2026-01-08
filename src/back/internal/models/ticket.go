package models

type Ticket struct {
	Number    int        `json:"number"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Number   int    `json:"number"`
	Section  string `json:"section"`
	Question string `json:"question"`
}

type TicketGenerationRequest struct {
	QuestionsPerTicket int `json:"questionsPerTicket" binding:"required,min=1,max=50"`
	TicketCount        int `json:"ticketCount" binding:"required,min=1,max=100"`
}
