# SEO Tech Platform

An AI-powered technical SEO and QA platform that analyzes websites, generates actionable SEO improvements, automation tests, and visual reports — designed for developers and engineering teams.

## 🏗️ Architecture

- **crawler-service** (Golang): High-performance web crawler using Colly
- **analyzer-service** (Python): AI-powered analysis using Lighthouse, Axe-core, and LangChain
- **api-gateway** (Golang): REST API backend and business logic
- **web-dashboard** (Vue.js): User interface for visualization and reporting

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+
- Python 3.11+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+

### Local Development

```bash
# Start all services
docker-compose up -d

# Access the dashboard
open http://localhost:3000

# API documentation
open http://localhost:8080/swagger
```

## 📦 Services

### Crawler Service (Port 8081)

- High-performance web crawling
- Headless browser support
- Queue-based job processing

### Analyzer Service (Port 8082)

- Lighthouse CI integration
- AI-powered suggestions
- Core Web Vitals analysis

### API Gateway (Port 8080)

- RESTful API
- User & Project management
- Database operations

### Web Dashboard (Port 3000)

- Real-time analytics
- Interactive reports
- Project management UI

## 🗄️ Database Schema

See `infrastructure/postgres/init.sql` for complete schema.

## 🛠️ Development

### Crawler Service

```bash
cd crawler-service
go mod download
go run cmd/worker/main.go
```

### Analyzer Service

```bash
cd analyzer-service
pip install -r requirements.txt
python src/main.py
```

### API Gateway

```bash
cd api-gateway
go mod download
go run cmd/server/main.go
```

### Web Dashboard

```bash
cd web-dashboard
npm install
npm run dev
```

## 📝 Environment Variables

Copy `.env.example` to `.env` and configure:

- Database credentials
- Redis connection
- OpenAI API key
- AWS credentials (for S3)

## 🚢 Deployment

### Kubernetes (AWS EKS)

```bash
cd infrastructure/terraform
terraform init
terraform apply

cd ../k8s
kubectl apply -f .
```

## 📊 Features

- ✅ Website crawling and analysis
- ✅ Core Web Vitals monitoring
- ✅ AI-powered SEO suggestions
- ✅ Lighthouse CI integration
- ✅ Automated testing
- ✅ Visual reporting
- ✅ Multi-project support

## 🤝 Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## 📄 License

MIT License - see LICENSE file for details
