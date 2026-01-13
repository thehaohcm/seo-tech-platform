# Build stage
FROM python:3.11-slim

# Install Node.js for Lighthouse
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    wget \
    chromium \
    chromium-driver \
    && curl -fsSL https://deb.nodesource.com/setup_18.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g lighthouse \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy requirements
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY src/ ./src/

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV CHROME_BIN=/usr/bin/chromium
ENV CHROME_PATH=/usr/bin/chromium

EXPOSE 8082

# Run the analyzer
CMD ["python", "src/main.py"]
