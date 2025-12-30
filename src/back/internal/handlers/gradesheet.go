package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/CreateLab/laritmo/internal/repository"
	"github.com/gin-gonic/gin"
)

type GradeSheetHandler struct {
	repo   *repository.GradeSheetRepository
	logger *slog.Logger
}

func NewGradeSheetHandler(repo *repository.GradeSheetRepository, logger *slog.Logger) *GradeSheetHandler {
	return &GradeSheetHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetAll godoc
// @Summary      Get all grade sheets
// @Description  Get list of all grade sheets with optional course filter
// @Tags         grade-sheets
// @Produce      json
// @Param        course_id  query     int  false  "Course ID filter"
// @Success      200        {array}   models.GradeSheet
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /api/grade-sheets [get]
func (h *GradeSheetHandler) GetAll(c *gin.Context) {
	var courseID *int
	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		id, err := strconv.Atoi(courseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
			return
		}
		courseID = &id
	}

	sheets, err := h.repo.GetAll(courseID)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get grade sheets", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get grade sheets"})
		return
	}

	if sheets == nil {
		sheets = []models.GradeSheet{}
	}

	c.JSON(http.StatusOK, sheets)
}

// GetByID godoc
// @Summary      Get grade sheet by ID
// @Description  Get grade sheet details by ID
// @Tags         grade-sheets
// @Produce      json
// @Param        id   path      int  true  "Grade Sheet ID"
// @Success      200  {object}  models.GradeSheet
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/grade-sheets/{id} [get]
func (h *GradeSheetHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	sheet, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get grade sheet", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get grade sheet"})
		return
	}

	if sheet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grade sheet not found"})
		return
	}

	c.JSON(http.StatusOK, sheet)
}

type CreateGradeSheetRequest struct {
	CourseID    int    `json:"course_id" binding:"required"`
	SheetURL    string `json:"sheet_url" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateGradeSheetRequest struct {
	SheetURL    string `json:"sheet_url" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// Create godoc
// @Summary      Create grade sheet
// @Description  Create a new grade sheet (admin only)
// @Tags         admin-gradesheets
// @Accept       json
// @Produce      json
// @Param        gradesheet  body      CreateGradeSheetRequest  true  "Grade sheet data"
// @Success      201         {object}  models.GradeSheet
// @Failure      400         {object}  map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      403         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/admin/grade-sheets [post]
func (h *GradeSheetHandler) Create(c *gin.Context) {
	var req CreateGradeSheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	gradeSheet, err := h.repo.Create(req.CourseID, req.SheetURL, req.Description)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to create grade sheet", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create grade sheet"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Grade sheet created", "id", gradeSheet.ID)
	c.JSON(http.StatusCreated, gradeSheet)
}

// Update godoc
// @Summary      Update grade sheet
// @Description  Update grade sheet by ID (admin only)
// @Tags         admin-gradesheets
// @Accept       json
// @Produce      json
// @Param        id          path      int                     true  "Grade Sheet ID"
// @Param        gradesheet  body      UpdateGradeSheetRequest  true  "Grade sheet data"
// @Success      200         {object}  map[string]string
// @Failure      400         {object}  map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      403         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/admin/grade-sheets/{id} [put]
func (h *GradeSheetHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req UpdateGradeSheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err = h.repo.Update(id, req.SheetURL, req.Description)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to update grade sheet", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update grade sheet"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Grade sheet updated", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "Grade sheet updated"})
}

// Delete godoc
// @Summary      Delete grade sheet
// @Description  Delete grade sheet by ID (admin only)
// @Tags         admin-gradesheets
// @Produce      json
// @Param        id   path      int  true  "Grade Sheet ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/admin/grade-sheets/{id} [delete]
func (h *GradeSheetHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to delete grade sheet", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete grade sheet"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Grade sheet deleted", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "Grade sheet deleted"})
}
