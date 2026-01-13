package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/seo-tech-platform/crawler-service/internal/engine"
	"github.com/seo-tech-platform/crawler-service/internal/queue"
	"github.com/seo-tech-platform/crawler-service/pkg/config"
	"github.com/seo-tech-platform/crawler-service/pkg/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Starting Crawler Service Worker...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize queue connection
	queueClient, err := queue.NewRedisQueue(cfg.RedisURL)
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer queueClient.Close()

	// Initialize crawler engine
	crawler := engine.NewCrawler(cfg, logger)

	// Start listening for crawl jobs
	logger.Info("Worker is ready and listening for jobs...")
	if err := queueClient.Listen("crawl_queue", crawler.ProcessJob); err != nil {
		logger.Fatalf("Failed to listen to queue: %v", err)
	}

	// Keep the worker running
	select {}
}
