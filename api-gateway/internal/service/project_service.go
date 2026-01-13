package service

import (
	"github.com/seo-tech-platform/api-gateway/internal/models"
	"gorm.io/gorm"
)

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

// ListProjects returns all projects for a user
func (s *ProjectService) ListProjects(userID int) ([]models.Project, error) {
	var projects []models.Project
	if err := s.db.Where("user_id = ?", userID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(userID int, domain, name string, settings map[string]interface{}) (*models.Project, error) {
	project := models.Project{
		UserID:   userID,
		Domain:   domain,
		Name:     name,
		Settings: settings,
	}

	if err := s.db.Create(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

// GetProject returns a project by ID
func (s *ProjectService) GetProject(id string) (*models.Project, error) {
	var project models.Project
	if err := s.db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(id, domain, name string, settings map[string]interface{}) (*models.Project, error) {
	var project models.Project
	if err := s.db.First(&project, id).Error; err != nil {
		return nil, err
	}

	if domain != "" {
		project.Domain = domain
	}
	if name != "" {
		project.Name = name
	}
	if settings != nil {
		project.Settings = settings
	}

	if err := s.db.Save(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(id string) error {
	return s.db.Delete(&models.Project{}, id).Error
}
