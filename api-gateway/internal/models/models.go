package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// JSONB is a custom type for handling JSONB fields
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var result interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	// Convert to map[string]interface{} if it's a map, otherwise keep as is (array, etc)
	switch v := result.(type) {
	case map[string]interface{}:
		*j = v
	case []interface{}:
		// For arrays, store as-is in a wrapper map with "items" key
		*j = map[string]interface{}{"items": v}
	default:
		*j = map[string]interface{}{"value": v}
	}

	return nil
}

// StringArray is a custom type for handling string array JSONB fields
type StringArray []string

// Value implements the driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal StringArray value: %v", value)
	}

	var result []string
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*s = result
	return nil
}

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
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"not null"`
	Domain    string    `json:"domain" gorm:"not null"`
	Name      string    `json:"name"`
	Settings  JSONB     `json:"settings" gorm:"type:jsonb"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	ID                  int         `json:"id" gorm:"primaryKey"`
	RunID               int         `json:"run_id" gorm:"not null"`
	URL                 string      `json:"url" gorm:"not null"`
	StatusCode          int         `json:"status_code"`
	LoadTimeMs          int         `json:"load_time_ms"`
	LCPScore            float64     `json:"lcp_score"`
	FIDScore            float64     `json:"fid_score"`
	CLSScore            float64     `json:"cls_score"`
	FCPScore            float64     `json:"fcp_score"`
	TTFBScore           float64     `json:"ttfb_score"`
	Title               string      `json:"title"`
	MetaDescription     string      `json:"meta_description"`
	H1Tags              StringArray `json:"h1_tags" gorm:"type:jsonb"`
	CanonicalURL        string      `json:"canonical_url"`
	HasRobotsMeta       bool        `json:"has_robots_meta"`
	SEOIssues           JSONB       `json:"seo_issues" gorm:"type:jsonb"`
	AccessibilityIssues JSONB       `json:"accessibility_issues" gorm:"type:jsonb"`
	PerformanceIssues   JSONB       `json:"performance_issues" gorm:"type:jsonb"`
	AISuggestions       string      `json:"ai_suggestions"`
	HTMLSnapshotPath    string      `json:"html_snapshot_path"`
	ScreenshotPath      string      `json:"screenshot_path"`
	CreatedAt           time.Time   `json:"created_at"`
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

	var db *gorm.DB
	var err error

	// Retry connection up to 30 times with 2 second delay
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Test the connection
			sqlDB, sqlErr := db.DB()
			if sqlErr == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					fmt.Printf("Successfully connected to database after %d attempts\n", i+1)
					return db, nil
				}
			}
		}

		fmt.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}
