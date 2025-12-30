package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/CreateLab/laritmo/internal/repository"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	repo   *repository.CourseRepository
	logger *slog.Logger
}

type CreateCourseRequest struct {
	Name        string `json:"name" binding:"required"`
	Semester    string `json:"semester" binding:"required"`
	Description string `json:"description"`
}

type UpdateCourseRequest struct {
	Name        string `json:"name" binding:"required"`
	Semester    string `json:"semester" binding:"required"`
	Description string `json:"description"`
}

func NewCourseHandler(repo *repository.CourseRepository, logger *slog.Logger) *CourseHandler {
	return &CourseHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetAll godoc
// @Summary      Get all courses
// @Description  Get list of all courses
// @Tags         courses
// @Produce      json
// @Success      200  {array}   models.Course
// @Failure      500  {object}  map[string]string
// @Router       /api/courses [get]
func (h *CourseHandler) GetAll(c *gin.Context) {
	courses, err := h.repo.GetAll()
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get courses", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	if courses == nil {
		courses = []models.Course{}
	}

	c.JSON(http.StatusOK, courses)
}

// GetByID godoc
// @Summary      Get course by ID
// @Description  Get course details by ID
// @Tags         courses
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  models.Course
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/courses/{id} [get]
func (h *CourseHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	course, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get course", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course"})
		return
	}

	if course == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// Create godoc
// @Summary      Create course
// @Description  Create a new course
// @Tags         admin-courses
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        course  body      CreateCourseRequest  true  "Course data"
// @Success      201     {object}  models.Course
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api/admin/courses [post]
func (h *CourseHandler) Create(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	course, err := h.repo.Create(req.Name, req.Semester, req.Description)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to create course", "error", err)
		c.JSON(500, gin.H{"error": "Failed to create course"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Course created", "id", course.ID)
	c.JSON(201, course)
}

// Update godoc
// @Summary      Update course
// @Description  Update course by ID
// @Tags         admin-courses
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path      int                  true  "Course ID"
// @Param        course  body      UpdateCourseRequest  true  "Course data"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api/admin/courses/{id} [put]
func (h *CourseHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.repo.Update(id, req.Name, req.Semester, req.Description)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to update course", "error", err)
		c.JSON(500, gin.H{"error": "Failed to update course"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Course updated", "id", id)
	c.JSON(200, gin.H{"message": "Course updated"})
}

// Delete godoc
// @Summary      Delete course
// @Description  Delete course by ID
// @Tags         admin-courses
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/admin/courses/{id} [delete]
func (h *CourseHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.repo.Delete(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to delete course", "error", err)
		c.JSON(500, gin.H{"error": "Failed to delete course"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Course deleted", "id", id)
	c.JSON(200, gin.H{"message": "Course deleted"})
}
