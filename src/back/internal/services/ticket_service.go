package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/CreateLab/laritmo/internal/models"
)

// ExamQuestionRepositoryInterface - interface for working with exam questions
type ExamQuestionRepositoryInterface interface {
	GetByCourseID(courseID int) ([]models.ExamQuestion, error)
}

type TicketService struct {
	examRepo ExamQuestionRepositoryInterface
}

func NewTicketService(examRepo ExamQuestionRepositoryInterface) *TicketService {
	return &TicketService{
		examRepo: examRepo,
	}
}

// GenerateRandomTicket generates a single random ticket from course questions
func (s *TicketService) GenerateRandomTicket(ctx context.Context, courseID int, questionsCount int) (*models.Ticket, error) {
	if questionsCount < 1 || questionsCount > 50 {
		return nil, errors.New("questions count must be between 1 and 50")
	}

	// Get all course questions
	allQuestions, err := s.examRepo.GetByCourseID(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam questions: %w", err)
	}

	if len(allQuestions) < questionsCount {
		return nil, fmt.Errorf("not enough questions: have %d, need %d", len(allQuestions), questionsCount)
	}

	// Group questions by section
	questionsBySection := make(map[string][]models.ExamQuestion)
	for _, q := range allQuestions {
		questionsBySection[q.Section] = append(questionsBySection[q.Section], q)
	}

	// Select questions according to the algorithm
	selectedQuestions := s.selectQuestions(questionsBySection, questionsCount)

	// Convert to Question format
	questions := make([]models.Question, len(selectedQuestions))
	for i, q := range selectedQuestions {
		questions[i] = models.Question{
			Number:   q.Number,
			Section:  q.Section,
			Question: q.Question,
		}
	}

	return &models.Ticket{
		Number:    1,
		Questions: questions,
	}, nil
}

// GenerateMultipleTickets generates multiple tickets with minimized question overlap
func (s *TicketService) GenerateMultipleTickets(ctx context.Context, courseID int, ticketCount, questionsPerTicket int) ([]models.Ticket, error) {
	if ticketCount < 1 || ticketCount > 100 {
		return nil, errors.New("ticket count must be between 1 and 100")
	}
	if questionsPerTicket < 1 || questionsPerTicket > 50 {
		return nil, errors.New("questions per ticket must be between 1 and 50")
	}

	// Get all course questions
	allQuestions, err := s.examRepo.GetByCourseID(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam questions: %w", err)
	}

	// Check if there are enough questions
	if len(allQuestions) < questionsPerTicket {
		return nil, fmt.Errorf("not enough questions: have %d, need at least %d", len(allQuestions), questionsPerTicket)
	}

	// Group questions by section
	questionsBySection := make(map[string][]models.ExamQuestion)
	for _, q := range allQuestions {
		questionsBySection[q.Section] = append(questionsBySection[q.Section], q)
	}

	// Generate tickets with tracking of used questions
	tickets := make([]models.Ticket, ticketCount)
	usedQuestions := make(map[int]int) // question ID -> count of usage

	for i := 0; i < ticketCount; i++ {
		selectedQuestions := s.selectQuestionsWithTracking(questionsBySection, questionsPerTicket, usedQuestions, allQuestions)

		questions := make([]models.Question, len(selectedQuestions))
		for j, q := range selectedQuestions {
			questions[j] = models.Question{
				Number:   q.Number,
				Section:  q.Section,
				Question: q.Question,
			}
			usedQuestions[q.ID]++
		}

		tickets[i] = models.Ticket{
			Number:    i + 1,
			Questions: questions,
		}
	}

	return tickets, nil
}

// selectQuestions selects questions according to the distribution algorithm
func (s *TicketService) selectQuestions(questionsBySection map[string][]models.ExamQuestion, questionsCount int) []models.ExamQuestion {
	var selected []models.ExamQuestion
	sections := make([]string, 0, len(questionsBySection))
	for section := range questionsBySection {
		sections = append(sections, section)
	}

	// If questions count >= sections count, take one from each section
	if questionsCount >= len(sections) {
		// Take one question from each section
		for _, section := range sections {
			sectionQuestions := questionsBySection[section]
			if len(sectionQuestions) > 0 {
				randomIndex := rand.Intn(len(sectionQuestions))
				selected = append(selected, sectionQuestions[randomIndex])
			}
		}

		// Fill the rest with random questions
		remaining := questionsCount - len(selected)
		if remaining > 0 {
			allQuestions := s.flattenQuestions(questionsBySection)
			selected = append(selected, s.selectRandomQuestions(allQuestions, remaining, selected)...)
		}
	} else {
		// Select random sections
		selectedSections := s.selectRandomSections(sections, questionsCount)
		for _, section := range selectedSections {
			sectionQuestions := questionsBySection[section]
			if len(sectionQuestions) > 0 {
				randomIndex := rand.Intn(len(sectionQuestions))
				selected = append(selected, sectionQuestions[randomIndex])
			}
		}
	}

	// Shuffle question order
	s.shuffleQuestions(selected)

	return selected
}

// selectQuestionsWithTracking selects questions considering already used ones
func (s *TicketService) selectQuestionsWithTracking(
	questionsBySection map[string][]models.ExamQuestion,
	questionsCount int,
	usedQuestions map[int]int,
	allQuestions []models.ExamQuestion,
) []models.ExamQuestion {
	var selected []models.ExamQuestion
	sections := make([]string, 0, len(questionsBySection))
	for section := range questionsBySection {
		sections = append(sections, section)
	}

	// Create list of available questions (priority to less used ones)
	availableQuestions := s.getAvailableQuestions(questionsBySection, usedQuestions)

	if questionsCount >= len(sections) {
		// Take one from each section (priority to least used)
		for _, section := range sections {
			sectionQuestions := questionsBySection[section]
			bestQuestion := s.findLeastUsedQuestion(sectionQuestions, usedQuestions)
			if bestQuestion != nil {
				selected = append(selected, *bestQuestion)
			}
		}

		// Fill the rest with least used questions
		remaining := questionsCount - len(selected)
		if remaining > 0 {
			additional := s.selectLeastUsedQuestions(availableQuestions, remaining, selected, usedQuestions)
			selected = append(selected, additional...)
		}
	} else {
		// Select random sections, but within them take least used questions
		selectedSections := s.selectRandomSections(sections, questionsCount)
		for _, section := range selectedSections {
			sectionQuestions := questionsBySection[section]
			bestQuestion := s.findLeastUsedQuestion(sectionQuestions, usedQuestions)
			if bestQuestion != nil {
				selected = append(selected, *bestQuestion)
			}
		}
	}

	// Shuffle question order
	s.shuffleQuestions(selected)

	return selected
}

// getAvailableQuestions returns all available questions considering used ones
func (s *TicketService) getAvailableQuestions(
	questionsBySection map[string][]models.ExamQuestion,
	usedQuestions map[int]int,
) []models.ExamQuestion {
	var all []models.ExamQuestion
	for _, questions := range questionsBySection {
		all = append(all, questions...)
	}
	return all
}

// findLeastUsedQuestion finds the least used question in a section
func (s *TicketService) findLeastUsedQuestion(questions []models.ExamQuestion, usedQuestions map[int]int) *models.ExamQuestion {
	if len(questions) == 0 {
		return nil
	}

	bestQuestion := &questions[0]
	bestCount := usedQuestions[questions[0].ID]

	for i := 1; i < len(questions); i++ {
		count := usedQuestions[questions[i].ID]
		if count < bestCount {
			bestCount = count
			bestQuestion = &questions[i]
		}
	}

	return bestQuestion
}

// selectLeastUsedQuestions selects the least used questions
func (s *TicketService) selectLeastUsedQuestions(
	availableQuestions []models.ExamQuestion,
	count int,
	alreadySelected []models.ExamQuestion,
	usedQuestions map[int]int,
) []models.ExamQuestion {
	// Create set of already selected IDs
	selectedIDs := make(map[int]bool)
	for _, q := range alreadySelected {
		selectedIDs[q.ID] = true
	}

	// Filter available questions (exclude already selected)
	filtered := make([]models.ExamQuestion, 0)
	for _, q := range availableQuestions {
		if !selectedIDs[q.ID] {
			filtered = append(filtered, q)
		}
	}

	// Sort by usage count (ascending)
	sorted := make([]models.ExamQuestion, len(filtered))
	copy(sorted, filtered)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if usedQuestions[sorted[i].ID] > usedQuestions[sorted[j].ID] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Take first count questions
	resultCount := count
	if resultCount > len(sorted) {
		resultCount = len(sorted)
	}

	return sorted[:resultCount]
}

// selectRandomSections selects random sections
func (s *TicketService) selectRandomSections(sections []string, count int) []string {
	if count >= len(sections) {
		return sections
	}

	selected := make([]string, 0, count)
	indices := rand.Perm(len(sections))
	for i := 0; i < count; i++ {
		selected = append(selected, sections[indices[i]])
	}

	return selected
}

// selectRandomQuestions selects random questions, excluding already selected
func (s *TicketService) selectRandomQuestions(allQuestions []models.ExamQuestion, count int, exclude []models.ExamQuestion) []models.ExamQuestion {
	excludeIDs := make(map[int]bool)
	for _, q := range exclude {
		excludeIDs[q.ID] = true
	}

	available := make([]models.ExamQuestion, 0)
	for _, q := range allQuestions {
		if !excludeIDs[q.ID] {
			available = append(available, q)
		}
	}

	if count > len(available) {
		count = len(available)
	}

	selected := make([]models.ExamQuestion, 0, count)
	indices := rand.Perm(len(available))
	for i := 0; i < count; i++ {
		selected = append(selected, available[indices[i]])
	}

	return selected
}

// flattenQuestions converts map to flat list
func (s *TicketService) flattenQuestions(questionsBySection map[string][]models.ExamQuestion) []models.ExamQuestion {
	var all []models.ExamQuestion
	for _, questions := range questionsBySection {
		all = append(all, questions...)
	}
	return all
}

// shuffleQuestions randomly shuffles questions
func (s *TicketService) shuffleQuestions(questions []models.ExamQuestion) {
	for i := len(questions) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		questions[i], questions[j] = questions[j], questions[i]
	}
}
