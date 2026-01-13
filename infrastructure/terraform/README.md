# Terraform Configuration for AWS Infrastructure

This directory contains Terraform configurations for provisioning AWS infrastructure.

## Resources

- VPC and networking
- RDS PostgreSQL database
- ElastiCache Redis
- EKS cluster
- S3 bucket for storage
- IAM roles and policies

## Usage

```bash
cd infrastructure/terraform

# Initialize Terraform
terraform init

# Plan infrastructure
terraform plan

# Apply infrastructure
terraform apply

# Destroy infrastructure (careful!)
terraform destroy
```

## Variables

Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your values:

```hcl
aws_region = "us-east-1"
project_name = "seo-platform"
environment = "production"
```

## Outputs

After applying, Terraform will output:

- RDS endpoint
- Redis endpoint
- EKS cluster name
- S3 bucket name
