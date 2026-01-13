package models

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	FullName     string    `json:"full_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Project represents a website project
type Project struct {
	ID        int                    `json:"id" gorm:"primaryKey"`
	UserID    int                    `json:"user_id" gorm:"not null"`
	Domain    string                 `json:"domain" gorm:"not null"`
	Name      string                 `json:"name"`
	Settings  map[string]interface{} `json:"settings" gorm:"type:jsonb"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// AuditRun represents an audit execution
type AuditRun struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	ProjectID    int        `json:"project_id" gorm:"not null"`
	Status       string     `json:"status" gorm:"default:'queued'"`
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	OverallScore *int       `json:"overall_score"`
	TotalPages   int        `json:"total_pages" gorm:"default:0"`
	ErrorMessage string     `json:"error_message,omitempty"`
}

// PageAudit represents audit results for a single page
type PageAudit struct {
	ID                  int                    `json:"id" gorm:"primaryKey"`
	RunID               int                    `json:"run_id" gorm:"not null"`
	URL                 string                 `json:"url" gorm:"not null"`
	StatusCode          int                    `json:"status_code"`
	LoadTimeMs          int                    `json:"load_time_ms"`
	LCPScore            float64                `json:"lcp_score"`
	FIDScore            float64                `json:"fid_score"`
	CLSScore            float64                `json:"cls_score"`
	FCPScore            float64                `json:"fcp_score"`
	TTFBScore           float64                `json:"ttfb_score"`
	Title               string                 `json:"title"`
	MetaDescription     string                 `json:"meta_description"`
	H1Tags              []string               `json:"h1_tags" gorm:"type:jsonb"`
	CanonicalURL        string                 `json:"canonical_url"`
	HasRobotsMeta       bool                   `json:"has_robots_meta"`
	SEOIssues           map[string]interface{} `json:"seo_issues" gorm:"type:jsonb"`
	AccessibilityIssues map[string]interface{} `json:"accessibility_issues" gorm:"type:jsonb"`
	PerformanceIssues   map[string]interface{} `json:"performance_issues" gorm:"type:jsonb"`
	AISuggestions       string                 `json:"ai_suggestions"`
	HTMLSnapshotPath    string                 `json:"html_snapshot_path"`
	ScreenshotPath      string                 `json:"screenshot_path"`
	CreatedAt           time.Time              `json:"created_at"`
}

// InitDB initializes database connection
func InitDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
