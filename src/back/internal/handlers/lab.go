package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/CreateLab/laritmo/internal/repository"
	"github.com/gin-gonic/gin"
)

type LabHandler struct {
	repo   *repository.LabRepository
	logger *slog.Logger
}

func NewLabHandler(repo *repository.LabRepository, logger *slog.Logger) *LabHandler {
	return &LabHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetAll godoc
// @Summary      Get all labs
// @Description  Get list of all labs with optional course filter
// @Tags         labs
// @Produce      json
// @Param        course_id  query     int  false  "Course ID filter"
// @Success      200        {array}   models.Lab
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /api/labs [get]
func (h *LabHandler) GetAll(c *gin.Context) {
	var courseID *int
	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		id, err := strconv.Atoi(courseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
			return
		}
		courseID = &id
	}

	labs, err := h.repo.GetAll(courseID)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get labs", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get labs"})
		return
	}

	if labs == nil {
		labs = []models.Lab{}
	}

	c.JSON(http.StatusOK, labs)
}

// GetByID godoc
// @Summary      Get lab by ID
// @Description  Get lab details by ID
// @Tags         labs
// @Produce      json
// @Param        id   path      int  true  "Lab ID"
// @Success      200  {object}  models.Lab
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/labs/{id} [get]
func (h *LabHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	lab, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get lab", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lab"})
		return
	}

	if lab == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lab not found"})
		return
	}

	c.JSON(http.StatusOK, lab)
}

type CreateLabRequest struct {
	CourseID    int     `json:"course_id" binding:"required"`
	Number      int     `json:"number" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	MaxScore    int     `json:"max_score" binding:"required"`
	GithubURL   string  `json:"github_url"`
	Deadline    *string `json:"deadline"`
}

type UpdateLabRequest struct {
	CourseID    int     `json:"course_id" binding:"required"`
	Number      int     `json:"number" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	MaxScore    int     `json:"max_score" binding:"required"`
	GithubURL   string  `json:"github_url"`
	Deadline    *string `json:"deadline"`
}

// Create godoc
// @Summary      Create lab
// @Description  Create a new lab
// @Tags         admin-labs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        lab   body      CreateLabRequest  true  "Lab data"
// @Success      201   {object}  models.Lab
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/admin/labs [post]
func (h *LabHandler) Create(c *gin.Context) {
	var req CreateLabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	lab, err := h.repo.Create(req.CourseID, req.Number, req.MaxScore, req.Title, req.Description, req.GithubURL, req.Deadline)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to create lab", "error", err)
		c.JSON(500, gin.H{"error": "Failed to create lab"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Lab created", "id", lab.ID)
	c.JSON(201, lab)
}

// Update godoc
// @Summary      Update lab
// @Description  Update lab by ID
// @Tags         admin-labs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Lab ID"
// @Param        lab  body      UpdateLabRequest   true  "Lab data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/admin/labs/{id} [put]
func (h *LabHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateLabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.repo.Update(id, req.CourseID, req.Number, req.MaxScore, req.Title, req.Description, req.GithubURL, req.Deadline)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to update lab", "error", err)
		c.JSON(500, gin.H{"error": "Failed to update lab"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Lab updated", "id", id)
	c.JSON(200, gin.H{"message": "Lab updated"})
}

// Delete godoc
// @Summary      Delete lab
// @Description  Delete lab by ID
// @Tags         admin-labs
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Lab ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/admin/labs/{id} [delete]
func (h *LabHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.repo.Delete(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to delete lab", "error", err)
		c.JSON(500, gin.H{"error": "Failed to delete lab"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Lab deleted", "id", id)
	c.JSON(200, gin.H{"message": "Lab deleted"})
}
