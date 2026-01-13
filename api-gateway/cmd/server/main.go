package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/seo-tech-platform/api-gateway/internal/api"
	"github.com/seo-tech-platform/api-gateway/internal/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := models.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate tables
	if err := db.AutoMigrate(&models.User{}, &models.Project{}, &models.AuditRun{}, &models.PageAudit{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Initialize API handlers
	apiHandler := api.NewHandler(db)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Authentication
		auth := v1.Group("/auth")
		{
			auth.POST("/register", apiHandler.Register)
			auth.POST("/login", apiHandler.Login)
		}

		// Projects (protected)
		projects := v1.Group("/projects")
		projects.Use(api.AuthMiddleware())
		{
			projects.GET("", apiHandler.ListProjects)
			projects.POST("", apiHandler.CreateProject)
			projects.GET("/:id", apiHandler.GetProject)
			projects.PUT("/:id", apiHandler.UpdateProject)
			projects.DELETE("/:id", apiHandler.DeleteProject)
		}

		// Audit Runs
		audits := v1.Group("/audits")
		audits.Use(api.AuthMiddleware())
		{
			audits.POST("/start", apiHandler.StartAudit)
			audits.GET("/:id", apiHandler.GetAuditRun)
			audits.GET("/project/:project_id", apiHandler.ListAuditRuns)
			audits.GET("/:id/pages", apiHandler.GetAuditPages)
		}
	}

	// Start server
	log.Println("Starting API Gateway on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
