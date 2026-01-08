package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/CreateLab/laritmo/internal/repository"
	"github.com/gin-gonic/gin"
)

type ExamQuestionHandler struct {
	repo   *repository.ExamQuestionRepository
	logger *slog.Logger
}

func NewExamQuestionHandler(repo *repository.ExamQuestionRepository, logger *slog.Logger) *ExamQuestionHandler {
	return &ExamQuestionHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetAll godoc
// @Summary      Get all exam questions
// @Description  Get list of all exam questions with optional course filter
// @Tags         exam-questions
// @Produce      json
// @Param        course_id  query     int  false  "Course ID filter"
// @Success      200        {array}   models.ExamQuestion
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /api/exam-questions [get]
func (h *ExamQuestionHandler) GetAll(c *gin.Context) {
	var questions []models.ExamQuestion
	var err error

	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
			return
		}
		questions, err = h.repo.GetByCourseID(courseID)
	} else {
		questions, err = h.repo.GetAll()
	}

	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get exam questions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exam questions"})
		return
	}

	if questions == nil {
		questions = []models.ExamQuestion{}
	}

	c.JSON(http.StatusOK, questions)
}

// GetByID godoc
// @Summary      Get exam question by ID
// @Description  Get exam question details by ID
// @Tags         exam-questions
// @Produce      json
// @Param        id   path      int  true  "Exam Question ID"
// @Success      200  {object}  models.ExamQuestion
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/exam-questions/{id} [get]
func (h *ExamQuestionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	question, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get exam question", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exam question"})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}

type CreateExamQuestionRequest struct {
	CourseID int    `json:"course_id" binding:"required"`
	Number   int    `json:"number" binding:"required"`
	Section  string `json:"section" binding:"required"`
	Question string `json:"question" binding:"required"`
}

type UpdateExamQuestionRequest struct {
	Number   int    `json:"number" binding:"required"`
	Section  string `json:"section" binding:"required"`
	Question string `json:"question" binding:"required"`
}

type BulkCreateJSONRequest struct {
	CourseID  int `json:"course_id" binding:"required"`
	Questions []struct {
		Number   int    `json:"number" binding:"required"`
		Section  string `json:"section" binding:"required"`
		Question string `json:"question" binding:"required"`
	} `json:"questions" binding:"required"`
}

// Create godoc
// @Summary      Create exam question
// @Description  Create a new exam question (admin only)
// @Tags         admin-exam-questions
// @Accept       json
// @Produce      json
// @Param        question  body      CreateExamQuestionRequest  true  "Exam question data"
// @Success      201       {object}  models.ExamQuestion
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      403       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/admin/exam-questions [post]
func (h *ExamQuestionHandler) Create(c *gin.Context) {
	var req CreateExamQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	question, err := h.repo.Create(req.CourseID, req.Number, req.Section, req.Question)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to create exam question", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exam question"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Exam question created", "id", question.ID)
	c.JSON(http.StatusCreated, question)
}

// Update godoc
// @Summary      Update exam question
// @Description  Update exam question by ID (admin only)
// @Tags         admin-exam-questions
// @Accept       json
// @Produce      json
// @Param        id         path      int                      true  "Exam Question ID"
// @Param        question   body      UpdateExamQuestionRequest  true  "Exam question data"
// @Success      200        {object}  map[string]string
// @Failure      400        {object}  map[string]string
// @Failure      401        {object}  map[string]string
// @Failure      403        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/admin/exam-questions/{id} [put]
func (h *ExamQuestionHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req UpdateExamQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	existingQuestion, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get exam question", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exam question"})
		return
	}
	if existingQuestion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam question not found"})
		return
	}

	err = h.repo.Update(id, existingQuestion.CourseID, req.Number, req.Section, req.Question)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to update exam question", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exam question"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Exam question updated", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "Exam question updated"})
}

// Delete godoc
// @Summary      Delete exam question
// @Description  Delete exam question by ID (admin only)
// @Tags         admin-exam-questions
// @Produce      json
// @Param        id   path      int  true  "Exam Question ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/admin/exam-questions/{id} [delete]
func (h *ExamQuestionHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to delete exam question", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exam question"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Exam question deleted", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "Exam question deleted"})
}

// BulkCreateJSON godoc
// @Summary      Bulk create exam questions from JSON
// @Description  Create multiple exam questions from JSON payload (admin only)
// @Tags         admin-exam-questions
// @Accept       json
// @Produce      json
// @Param        data   body      BulkCreateJSONRequest  true  "Bulk exam questions data"
// @Success      201    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/admin/exam-questions/bulk [post]
func (h *ExamQuestionHandler) BulkCreateJSON(c *gin.Context) {
	var req BulkCreateJSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(req.Questions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Questions list is empty"})
		return
	}

	var questions []models.ExamQuestion
	for _, q := range req.Questions {
		questions = append(questions, models.ExamQuestion{
			CourseID: req.CourseID,
			Number:   q.Number,
			Section:  q.Section,
			Question: q.Question,
		})
	}

	err := h.repo.BulkCreate(questions)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Bulk creation error exam questions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bulk creation error exam questions"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Exam questions bulk created", "count", len(questions), "course_id", req.CourseID)
	c.JSON(http.StatusCreated, gin.H{"message": "Questions successfully created", "count": len(questions)})
}

// BulkUploadFile godoc
// @Summary      Bulk upload exam questions from file
// @Description  Upload exam questions from JSON or CSV file (admin only)
// @Tags         admin-exam-questions
// @Accept       multipart/form-data
// @Produce      json
// @Param        course_id  formData  int     true   "Course ID"
// @Param        file       formData  file    true   "JSON or CSV file"
// @Success      201        {object}  map[string]string
// @Failure      400        {object}  map[string]string
// @Failure      401        {object}  map[string]string
// @Failure      403        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/admin/exam-questions/upload [post]
func (h *ExamQuestionHandler) BulkUploadFile(c *gin.Context) {
	courseIDStr := c.PostForm("course_id")
	if courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id is required"})
		return
	}

	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	var questions []models.ExamQuestion

	switch ext {
	case ".json":
		var jsonData map[string]interface{}
		if err := json.NewDecoder(file).Decode(&jsonData); err != nil {
			h.logger.ErrorContext(c.Request.Context(), "Parsing error JSON", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON format: %v", err)})
			return
		}

		questionsRaw, ok := jsonData["questions"]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON must contain 'questions' field"})
			return
		}

		questionsArray, ok := questionsRaw.([]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "'questions' field must be an array"})
			return
		}

		for _, qRaw := range questionsArray {
			qMap, ok := qRaw.(map[string]interface{})
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Each 'questions' element must be an object"})
				return
			}

			number, ok := qMap["number"].(float64)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "'number' field is required and must be a number"})
				return
			}

			section, ok := qMap["section"].(string)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "'section' field is required and must be a string"})
				return
			}

			question, ok := qMap["question"].(string)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "'question' field is required and must be a string"})
				return
			}

			questions = append(questions, models.ExamQuestion{
				CourseID: courseID,
				Number:   int(number),
				Section:  section,
				Question: question,
			})
		}

	case ".csv":
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			h.logger.ErrorContext(c.Request.Context(), "Reading error CSV", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to read CSV file: %v", err)})
			return
		}

		if len(records) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file must contain header and at least one data row"})
			return
		}

		for i, record := range records[1:] {
			if len(record) < 3 {
				h.logger.ErrorContext(c.Request.Context(), "Invalid number of columns in CSV", "row", i+2, "columns", len(record))
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("CSV must contain columns: number,section,question. Row %d has %d columns", i+2, len(record))})
				return
			}

			number, err := strconv.Atoi(record[0])
			if err != nil {
				h.logger.ErrorContext(c.Request.Context(), "Failed to parse number in CSV", "row", i+2, "error", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid number format in row %d: %v", i+2, err)})
				return
			}

			questions = append(questions, models.ExamQuestion{
				CourseID: courseID,
				Number:   number,
				Section:  record[1],
				Question: record[2],
			})
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only .json and .csv files are supported"})
		return
	}

	if len(questions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No questions found for import"})
		return
	}

	err = h.repo.BulkCreate(questions)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Bulk creation error exam questions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create questions: %v", err)})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Exam questions loaded from file", "count", len(questions), "course_id", courseID, "format", ext)
	c.JSON(http.StatusCreated, gin.H{"message": "Questions successfully loaded", "count": len(questions)})
}
