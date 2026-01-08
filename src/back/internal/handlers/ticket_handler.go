package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/gin-gonic/gin"
)

// TicketServiceInterface - интерфейс для сервиса генерации билетов
type TicketServiceInterface interface {
	GenerateRandomTicket(ctx context.Context, courseID int, questionsCount int) (*models.Ticket, error)
	GenerateMultipleTickets(ctx context.Context, courseID int, ticketCount, questionsPerTicket int) ([]models.Ticket, error)
}

// DocumentServiceInterface - интерфейс для сервиса генерации документов
type DocumentServiceInterface interface {
	GenerateTicketsDocument(tickets []models.Ticket) []byte
}

// CourseRepositoryInterface - интерфейс для репозитория курсов
type CourseRepositoryInterface interface {
	GetByID(id int) (*models.Course, error)
}

type TicketHandler struct {
	ticketService   TicketServiceInterface
	documentService DocumentServiceInterface
	courseRepo      CourseRepositoryInterface
	logger          *slog.Logger
}

func NewTicketHandler(
	ticketService TicketServiceInterface,
	documentService DocumentServiceInterface,
	courseRepo CourseRepositoryInterface,
	logger *slog.Logger,
) *TicketHandler {
	return &TicketHandler{
		ticketService:   ticketService,
		documentService: documentService,
		courseRepo:      courseRepo,
		logger:          logger,
	}
}

// GetRandomTicket godoc
// @Summary      Get random ticket
// @Description  Generate and return a random exam ticket for a course
// @Tags         tickets
// @Produce      json
// @Param        id         path      int  true   "Course ID"
// @Param        questions  query     int  false  "Number of questions per ticket (1-50)" default(10)
// @Success      200        {object}  map[string]models.Ticket
// @Failure      400        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /api/courses/{id}/tickets/random [get]
func (h *TicketHandler) GetRandomTicket(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Invalid course ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Проверяем существование курса
	course, err := h.courseRepo.GetByID(courseID)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get course", "error", err, "course_id", courseID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course"})
		return
	}
	if course == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Получаем количество вопросов из query параметра
	questionsStr := c.DefaultQuery("questions", "10")
	questionsCount, err := strconv.Atoi(questionsStr)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Invalid questions parameter", "error", err, "questions", questionsStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid questions parameter"})
		return
	}

	// Валидация количества вопросов
	if questionsCount < 1 || questionsCount > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Questions count must be between 1 and 50"})
		return
	}

	// Генерируем билет
	ticket, err := h.ticketService.GenerateRandomTicket(c.Request.Context(), courseID, questionsCount)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to generate ticket", "error", err, "course_id", courseID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

// GenerateTicketsDocument godoc
// @Summary      Generate tickets document
// @Description  Generate multiple exam tickets and return as TXT file (admin only)
// @Tags         admin-tickets
// @Accept       json
// @Produce      text/plain
// @Param        id      path      int                          true  "Course ID"
// @Param        request body      models.TicketGenerationRequest  true  "Generation parameters"
// @Success      200     {file}    binary                      "TXT file with tickets"
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/courses/{id}/tickets/generate [post]
func (h *TicketHandler) GenerateTicketsDocument(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Invalid course ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Проверяем существование курса
	course, err := h.courseRepo.GetByID(courseID)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get course", "error", err, "course_id", courseID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course"})
		return
	}
	if course == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Валидация body
	var req models.TicketGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Генерируем билеты
	tickets, err := h.ticketService.GenerateMultipleTickets(c.Request.Context(), courseID, req.TicketCount, req.QuestionsPerTicket)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to generate tickets", "error", err, "course_id", courseID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tickets"})
		return
	}

	// Генерируем TXT документ
	document := h.documentService.GenerateTicketsDocument(tickets)

	// Создаем имя файла из названия курса
	courseSlug := strings.ToLower(strings.ReplaceAll(course.Name, " ", "_"))
	courseSlug = strings.ReplaceAll(courseSlug, "/", "_")
	filename := "tickets_" + courseSlug + ".txt"

	// Устанавливаем заголовки для скачивания файла
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/plain; charset=utf-8", document)

	h.logger.InfoContext(c.Request.Context(), "Tickets document generated", "course_id", courseID, "ticket_count", len(tickets))
}
