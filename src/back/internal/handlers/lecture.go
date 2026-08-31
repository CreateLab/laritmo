package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/CreateLab/laritmo/internal/repository"
	"github.com/gin-gonic/gin"
)

type LectureHandler struct {
	repo       *repository.LectureRepository
	courseRepo *repository.CourseRepository
	logger     *slog.Logger
}

func NewLectureHandler(repo *repository.LectureRepository, courseRepo *repository.CourseRepository, logger *slog.Logger) *LectureHandler {
	return &LectureHandler{
		repo:       repo,
		courseRepo: courseRepo,
		logger:     logger,
	}
}

// GetAll godoc
// @Summary      Get all lectures
// @Description  Get list of all lectures with optional course filter
// @Tags         lectures
// @Produce      json
// @Param        course_id  query     int  false  "Course ID filter"
// @Success      200        {array}   models.Lecture
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /api/lectures [get]
func (h *LectureHandler) GetAll(c *gin.Context) {
	var courseID *int
	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		id, err := strconv.Atoi(courseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
			return
		}
		courseID = &id
	}

	lectures, err := h.repo.GetAll(courseID)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get lectures", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lectures"})
		return
	}

	if lectures == nil {
		lectures = []models.Lecture{}
	}

	c.JSON(http.StatusOK, lectures)
}

// GetByID godoc
// @Summary      Get lecture by ID
// @Description  Get lecture details by ID
// @Tags         lectures
// @Produce      json
// @Param        id   path      int  true  "Lecture ID"
// @Success      200  {object}  models.Lecture
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/lectures/{id} [get]
func (h *LectureHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	lecture, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get lecture", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lecture"})
		return
	}

	if lecture == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lecture not found"})
		return
	}

	c.JSON(http.StatusOK, lecture)
}

// ExportMarkdown godoc
// @Summary      Export lectures as Markdown
// @Description  Download all lectures (optionally filtered by course) as a single Markdown file. Requires an authenticated user.
// @Tags         lectures
// @Security     BearerAuth
// @Produce      text/markdown
// @Param        course_id  query  int  false  "Course ID filter"
// @Success      200  {string}  string  "Markdown file"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/lectures/export [get]
func (h *LectureHandler) ExportMarkdown(c *gin.Context) {
	var courseID *int
	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		id, err := strconv.Atoi(courseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
			return
		}
		courseID = &id
	}

	lectures, err := h.repo.GetAll(courseID)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get lectures", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lectures"})
		return
	}

	courses, err := h.courseRepo.GetAll()
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to get courses", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	courseNames := make(map[int]string, len(courses))
	for _, course := range courses {
		courseNames[course.ID] = course.Name
	}

	var sb strings.Builder
	lastCourseID := -1
	for _, l := range lectures {
		if l.CourseID != lastCourseID {
			if lastCourseID != -1 {
				sb.WriteString("\n")
			}
			name := courseNames[l.CourseID]
			if name == "" {
				name = fmt.Sprintf("Course %d", l.CourseID)
			}
			sb.WriteString(fmt.Sprintf("# %s\n\n", name))
			lastCourseID = l.CourseID
		}

		sb.WriteString(fmt.Sprintf("## Week %d: %s\n\n", l.Week, l.Title))
		if l.GithubURL != nil && *l.GithubURL != "" {
			sb.WriteString(fmt.Sprintf("GitHub: %s\n\n", *l.GithubURL))
		}
		sb.WriteString(l.Content)
		sb.WriteString("\n\n")
	}

	filename := "lectures.md"
	if courseID != nil {
		filename = fmt.Sprintf("lectures-course-%d.md", *courseID)
	}

	userID, _ := c.Get("user_id")
	h.logger.InfoContext(c.Request.Context(), "Lectures exported", "user_id", userID, "course_id", courseID)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Data(http.StatusOK, "text/markdown; charset=utf-8", []byte(sb.String()))
}

type CreateLectureRequest struct {
	CourseID  int    `json:"course_id" binding:"required"`
	Week      int    `json:"week" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	GithubURL string `json:"github_url"`
}

type UpdateLectureRequest struct {
	CourseID  int    `json:"course_id" binding:"required"`
	Week      int    `json:"week" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	GithubURL string `json:"github_url"`
}

// Create godoc
// @Summary      Create lecture
// @Description  Create a new lecture
// @Tags         admin-lectures
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        lecture  body      CreateLectureRequest  true  "Lecture data"
// @Success      201      {object}  models.Lecture
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/admin/lectures [post]
func (h *LectureHandler) Create(c *gin.Context) {
	var req CreateLectureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	lecture, err := h.repo.Create(req.CourseID, req.Week, req.Title, req.Content, req.GithubURL)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to create lecture", "error", err)
		c.JSON(500, gin.H{"error": "Failed to create lecture"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Lecture created", "id", lecture.ID)
	c.JSON(201, lecture)
}

// Update godoc
// @Summary      Update lecture
// @Description  Update lecture by ID
// @Tags         admin-lectures
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                   true  "Lecture ID"
// @Param        lecture  body      UpdateLectureRequest  true  "Lecture data"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/admin/lectures/{id} [put]
func (h *LectureHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateLectureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Validation error", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.repo.Update(id, req.CourseID, req.Week, req.Title, req.Content, req.GithubURL)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to update lecture", "error", err)
		c.JSON(500, gin.H{"error": "Failed to update lecture"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Lecture updated", "id", id)
	c.JSON(200, gin.H{"message": "Lecture updated"})
}

// Delete godoc
// @Summary      Delete lecture
// @Description  Delete lecture by ID
// @Tags         admin-lectures
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Lecture ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/admin/lectures/{id} [delete]
func (h *LectureHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.repo.Delete(id)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "Failed to delete lecture", "error", err)
		c.JSON(500, gin.H{"error": "Failed to delete lecture"})
		return
	}

	h.logger.InfoContext(c.Request.Context(), "Lecture deleted", "id", id)
	c.JSON(200, gin.H{"message": "Lecture deleted"})
}
