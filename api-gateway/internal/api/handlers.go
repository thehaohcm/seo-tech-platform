package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seo-tech-platform/api-gateway/internal/models"
	"github.com/seo-tech-platform/api-gateway/internal/queue"
	"github.com/seo-tech-platform/api-gateway/internal/service"
	"gorm.io/gorm"
)

type Handler struct {
	db          *gorm.DB
	service     *service.ProjectService
	authService *service.AuthService
	queue       *queue.RedisQueue
}

func NewHandler(db *gorm.DB, redisQueue *queue.RedisQueue) *Handler {
	return &Handler{
		db:          db,
		service:     service.NewProjectService(db),
		authService: service.NewAuthService(db),
		queue:       redisQueue,
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles user authentication
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// ListProjects returns all projects for the authenticated user
func (h *Handler) ListProjects(c *gin.Context) {
	userID := c.GetInt("user_id")

	projects, err := h.service.ListProjects(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// CreateProject creates a new project
func (h *Handler) CreateProject(c *gin.Context) {
	var req struct {
		Domain   string       `json:"domain" binding:"required"`
		Name     string       `json:"name" binding:"required"`
		Settings models.JSONB `json:"settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")

	project, err := h.service.CreateProject(userID, req.Domain, req.Name, req.Settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"project": project})
}

// GetProject returns a specific project
func (h *Handler) GetProject(c *gin.Context) {
	id := c.Param("id")

	project, err := h.service.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// UpdateProject updates a project
func (h *Handler) UpdateProject(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Domain   string       `json:"domain"`
		Name     string       `json:"name"`
		Settings models.JSONB `json:"settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.service.UpdateProject(id, req.Domain, req.Name, req.Settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// DeleteProject deletes a project
func (h *Handler) DeleteProject(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// StartAudit initiates a new audit run
func (h *Handler) StartAudit(c *gin.Context) {
	var req struct {
		ProjectID int `json:"project_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get project details
	var project models.Project
	if err := h.db.First(&project, req.ProjectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Create audit run
	auditRun := models.AuditRun{
		ProjectID: req.ProjectID,
		Status:    "queued",
		StartedAt: time.Now(),
	}

	if err := h.db.Create(&auditRun).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Push job to crawler queue
	crawlJob := queue.CrawlJob{
		RunID:     auditRun.ID,
		ProjectID: project.ID,
		StartURL:  project.Domain,
		MaxPages:  100, // Default max pages
	}

	if err := h.queue.PushCrawlJob(crawlJob); err != nil {
		// Update audit run status to failed
		auditRun.Status = "failed"
		auditRun.ErrorMessage = "Failed to queue crawl job: " + err.Error()
		h.db.Save(&auditRun)
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue crawl job"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"audit_run": auditRun})
}

// GetAuditRun returns details of an audit run
func (h *Handler) GetAuditRun(c *gin.Context) {
	id := c.Param("id")

	var auditRun models.AuditRun
	if err := h.db.First(&auditRun, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audit run not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"audit_run": auditRun})
}

// ListAuditRuns returns all audit runs for a project
func (h *Handler) ListAuditRuns(c *gin.Context) {
	projectID := c.Param("project_id")

	var auditRuns []models.AuditRun
	if err := h.db.Where("project_id = ?", projectID).Order("started_at DESC").Find(&auditRuns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"audit_runs": auditRuns})
}

// GetAuditPages returns all page audits for an audit run
func (h *Handler) GetAuditPages(c *gin.Context) {
	runID := c.Param("id")

	var pages []models.PageAudit
	if err := h.db.Where("run_id = ?", runID).Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages})
}
