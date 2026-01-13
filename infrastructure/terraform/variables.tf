variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "seo-platform"
}

variable "environment" {
  description = "Environment (dev, staging, production)"
  type        = string
  default     = "production"
}

variable "db_username" {
  description = "Database username"
  type        = string
  default     = "seouser"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}
