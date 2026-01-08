package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/CreateLab/laritmo/internal/models"
)

// ExamQuestionRepositoryInterface - интерфейс для работы с экзаменационными вопросами
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

// GenerateRandomTicket генерирует один случайный билет из вопросов курса
func (s *TicketService) GenerateRandomTicket(ctx context.Context, courseID int, questionsCount int) (*models.Ticket, error) {
	if questionsCount < 1 || questionsCount > 50 {
		return nil, errors.New("questions count must be between 1 and 50")
	}

	// Получаем все вопросы курса
	allQuestions, err := s.examRepo.GetByCourseID(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam questions: %w", err)
	}

	if len(allQuestions) < questionsCount {
		return nil, fmt.Errorf("not enough questions: have %d, need %d", len(allQuestions), questionsCount)
	}

	// Группируем вопросы по разделам
	questionsBySection := make(map[string][]models.ExamQuestion)
	for _, q := range allQuestions {
		questionsBySection[q.Section] = append(questionsBySection[q.Section], q)
	}

	// Выбираем вопросы согласно алгоритму
	selectedQuestions := s.selectQuestions(questionsBySection, questionsCount)

	// Преобразуем в формат Question
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

// GenerateMultipleTickets генерирует несколько билетов с минимизацией пересечений
func (s *TicketService) GenerateMultipleTickets(ctx context.Context, courseID int, ticketCount, questionsPerTicket int) ([]models.Ticket, error) {
	if ticketCount < 1 || ticketCount > 100 {
		return nil, errors.New("ticket count must be between 1 and 100")
	}
	if questionsPerTicket < 1 || questionsPerTicket > 50 {
		return nil, errors.New("questions per ticket must be between 1 and 50")
	}

	// Получаем все вопросы курса
	allQuestions, err := s.examRepo.GetByCourseID(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam questions: %w", err)
	}

	// Проверяем достаточность вопросов
	if len(allQuestions) < questionsPerTicket {
		return nil, fmt.Errorf("not enough questions: have %d, need at least %d", len(allQuestions), questionsPerTicket)
	}

	// Группируем вопросы по разделам
	questionsBySection := make(map[string][]models.ExamQuestion)
	for _, q := range allQuestions {
		questionsBySection[q.Section] = append(questionsBySection[q.Section], q)
	}

	// Генерируем билеты с отслеживанием использованных вопросов
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

// selectQuestions выбирает вопросы согласно алгоритму распределения
func (s *TicketService) selectQuestions(questionsBySection map[string][]models.ExamQuestion, questionsCount int) []models.ExamQuestion {
	var selected []models.ExamQuestion
	sections := make([]string, 0, len(questionsBySection))
	for section := range questionsBySection {
		sections = append(sections, section)
	}

	// Если количество вопросов >= количества разделов, берем по одному из каждого раздела
	if questionsCount >= len(sections) {
		// Берем по одному вопросу из каждого раздела
		for _, section := range sections {
			sectionQuestions := questionsBySection[section]
			if len(sectionQuestions) > 0 {
				randomIndex := rand.Intn(len(sectionQuestions))
				selected = append(selected, sectionQuestions[randomIndex])
			}
		}

		// Остальное заполняем случайными вопросами
		remaining := questionsCount - len(selected)
		if remaining > 0 {
			allQuestions := s.flattenQuestions(questionsBySection)
			selected = append(selected, s.selectRandomQuestions(allQuestions, remaining, selected)...)
		}
	} else {
		// Выбираем случайные разделы
		selectedSections := s.selectRandomSections(sections, questionsCount)
		for _, section := range selectedSections {
			sectionQuestions := questionsBySection[section]
			if len(sectionQuestions) > 0 {
				randomIndex := rand.Intn(len(sectionQuestions))
				selected = append(selected, sectionQuestions[randomIndex])
			}
		}
	}

	// Перемешиваем порядок вопросов
	s.shuffleQuestions(selected)

	return selected
}

// selectQuestionsWithTracking выбирает вопросы с учетом уже использованных
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

	// Создаем список доступных вопросов (приоритет тем, которые использовались меньше)
	availableQuestions := s.getAvailableQuestions(questionsBySection, usedQuestions)

	if questionsCount >= len(sections) {
		// Берем по одному из каждого раздела (приоритет менее использованным)
		for _, section := range sections {
			sectionQuestions := questionsBySection[section]
			bestQuestion := s.findLeastUsedQuestion(sectionQuestions, usedQuestions)
			if bestQuestion != nil {
				selected = append(selected, *bestQuestion)
			}
		}

		// Остальное заполняем наименее использованными вопросами
		remaining := questionsCount - len(selected)
		if remaining > 0 {
			additional := s.selectLeastUsedQuestions(availableQuestions, remaining, selected, usedQuestions)
			selected = append(selected, additional...)
		}
	} else {
		// Выбираем случайные разделы, но внутри них берем наименее использованные вопросы
		selectedSections := s.selectRandomSections(sections, questionsCount)
		for _, section := range selectedSections {
			sectionQuestions := questionsBySection[section]
			bestQuestion := s.findLeastUsedQuestion(sectionQuestions, usedQuestions)
			if bestQuestion != nil {
				selected = append(selected, *bestQuestion)
			}
		}
	}

	// Перемешиваем порядок вопросов
	s.shuffleQuestions(selected)

	return selected
}

// getAvailableQuestions возвращает все доступные вопросы с учетом использованных
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

// findLeastUsedQuestion находит наименее использованный вопрос в секции
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

// selectLeastUsedQuestions выбирает наименее использованные вопросы
func (s *TicketService) selectLeastUsedQuestions(
	availableQuestions []models.ExamQuestion,
	count int,
	alreadySelected []models.ExamQuestion,
	usedQuestions map[int]int,
) []models.ExamQuestion {
	// Создаем множество уже выбранных ID
	selectedIDs := make(map[int]bool)
	for _, q := range alreadySelected {
		selectedIDs[q.ID] = true
	}

	// Фильтруем доступные вопросы (исключаем уже выбранные)
	filtered := make([]models.ExamQuestion, 0)
	for _, q := range availableQuestions {
		if !selectedIDs[q.ID] {
			filtered = append(filtered, q)
		}
	}

	// Сортируем по количеству использований (ascending)
	sorted := make([]models.ExamQuestion, len(filtered))
	copy(sorted, filtered)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if usedQuestions[sorted[i].ID] > usedQuestions[sorted[j].ID] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Берем первые count вопросов
	resultCount := count
	if resultCount > len(sorted) {
		resultCount = len(sorted)
	}

	return sorted[:resultCount]
}

// selectRandomSections выбирает случайные разделы
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

// selectRandomQuestions выбирает случайные вопросы, исключая уже выбранные
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

// flattenQuestions преобразует map в плоский список
func (s *TicketService) flattenQuestions(questionsBySection map[string][]models.ExamQuestion) []models.ExamQuestion {
	var all []models.ExamQuestion
	for _, questions := range questionsBySection {
		all = append(all, questions...)
	}
	return all
}

// shuffleQuestions перемешивает вопросы случайным образом
func (s *TicketService) shuffleQuestions(questions []models.ExamQuestion) {
	for i := len(questions) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		questions[i], questions[j] = questions[j], questions[i]
	}
}
