# SEO Tech Platform

An AI-powered technical SEO and QA platform that analyzes websites, generates actionable SEO improvements, automation tests with Selenium, and visual reports — designed for developers and engineering teams.

## 🏗️ Architecture

- **crawler-service** (Go 1.23): High-performance web crawler using Colly with domain filtering
- **analyzer-service** (Python 3.11): AI-powered analysis with Lighthouse, Selenium, and LangChain
- **api-gateway** (Go 1.23): REST API backend with JWT authentication
- **web-dashboard** (Vue.js 3): Real-time dashboard with Composition API and Tailwind CSS

## 🛠️ Tech Stack

### Backend
- **Go 1.23** - API Gateway & Crawler Service
  - Gin Web Framework
  - GORM (PostgreSQL ORM)
  - Colly (Web Scraping)
  - JWT Authentication
  
- **Python 3.11** - Analyzer Service
  - Selenium WebDriver (Automated Testing)
  - LangChain + OpenAI GPT-4 (AI Suggestions)
  - Lighthouse CI (Performance Audits)
  - SQLAlchemy (Database ORM)

### Frontend
- **Vue.js 3** with Composition API
- **Pinia** - State Management
- **Tailwind CSS** - Styling
- **Vite** - Build Tool
- **Axios** - HTTP Client

### Infrastructure
- **PostgreSQL 15** - Primary Database with JSONB support
- **Redis 7** - Queue Management (crawl_queue, analysis_queue, test_queue)
- **Docker & Docker Compose** - Containerization
- **Nginx** - Reverse Proxy

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.23+
- Python 3.11+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- OpenAI API Key (for AI suggestions)

### Local Development

```bash
# Start all services
docker-compose up -d

# Access the dashboard
open http://localhost:8080

# Check service health
docker ps

# View logs
docker logs -f seo-api
docker logs -f seo-analyzer
docker logs -f seo-crawler
```

## 📦 Services

### Crawler Service

- High-performance web crawling with Colly
- Domain filtering (only crawls same-domain URLs)
- Redis queue-based job processing
- Parallel crawling with rate limiting
- Extracts: titles, meta descriptions, H1 tags, links, status codes

### Analyzer Service

- **Lighthouse CI** - Performance, SEO, Accessibility audits
- **Selenium WebDriver** - Automated testing with Chrome headless
  - 8 automated test types (page load, title, meta, H1, images, forms, links, console errors)
  - Screenshot capture (base64 format)
- **LangChain + GPT-4** - AI-powered SEO improvement suggestions
- **Redis Results Storage** - Test results with 1-hour expiry
- Processes jobs from analysis_queue and test_queue

### API Gateway (Port 8080)

- RESTful API with Gin framework
- JWT Authentication & Authorization
- User & Project management
- CRUD operations for audits and pages
- Redis integration for queue management
- PostgreSQL with JSONB support

### Web Dashboard (Port 8080)

- Real-time audit results with auto-polling (2-3 second intervals)
- SEO rating badges (Good/Average/Poor)
- Animated loading states with progress indicators
- Per-page automated testing with "Generate Auto Test" button
- Screenshot viewer for test results
- Project management (create, delete projects)
- Interactive reports with detailed metrics

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

```bash
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/seo_platform

# Redis
REDIS_URL=redis://localhost:6379

# OpenAI API (for AI suggestions)
OPENAI_API_KEY=sk-...

# JWT Secret
JWT_SECRET=your-secret-key

# Ports
API_PORT=8080
DASHBOARD_PORT=8080
```

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

### Core Features
- ✅ **Smart Web Crawling** - Domain-filtered crawling with same-origin policy
- ✅ **Lighthouse Audits** - Performance, SEO, Accessibility, Best Practices
- ✅ **AI-Powered Suggestions** - GPT-4 generates actionable SEO improvements
- ✅ **Automated Testing** - Selenium-based tests for each page
  - Page load validation
  - Title & meta description checks
  - H1 tag analysis
  - Image alt text validation
  - Form label accessibility
  - Link validity checks
  - Console error detection
- ✅ **Screenshot Capture** - Visual page snapshots during testing
- ✅ **Real-time Updates** - Auto-polling dashboard with live results
- ✅ **SEO Ratings** - Good/Average/Poor badges based on scores
- ✅ **Multi-project Support** - Manage multiple websites
- ✅ **Project Management** - Create, view, delete projects

### Technical Features
- 🔐 **JWT Authentication** - Secure user sessions
- 📊 **JSONB Storage** - Flexible schema for audit data
- 🚀 **Redis Queues** - Distributed job processing
- 📸 **Base64 Screenshots** - Embedded image storage
- ⏱️ **Rate Limiting** - Controlled crawling speed
- 🎯 **Domain Filtering** - Prevents external URL crawling
- 💾 **Result Caching** - 1-hour Redis expiry for test results

## 🤝 Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## 📄 License

MIT License - see LICENSE file for details
