# Analyzer Service

AI-powered SEO and performance analysis service built with Python.

## Features

- 🔍 Lighthouse CI integration
- ♿ Accessibility audits with Axe-core
- 🤖 AI-powered suggestions using OpenAI
- 📊 Core Web Vitals analysis
- 🔄 Redis-based job processing

## Development

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Install Node.js dependencies (for Lighthouse)
npm install -g lighthouse

# Run locally
python src/main.py
```

## Configuration

Set environment variables in `.env`:

- `OPENAI_API_KEY`: OpenAI API key for AI suggestions
- `REDIS_URL`: Redis connection string
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`: Database configuration

## Architecture

```
analyzer-service/
├── src/
│   ├── audits/              - Lighthouse and accessibility audits
│   ├── ai_agent/            - AI suggestion generation
│   ├── models/              - Database models
│   ├── queue_listener.py    - Redis queue consumer
│   └── main.py              - Entry point
└── requirements.txt
```

## Requirements

- Python 3.11+
- Node.js 18+ (for Lighthouse)
- Chrome/Chromium (for headless testing)
