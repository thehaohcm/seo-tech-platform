package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisQueue(redisURL string) (*RedisQueue, error) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisQueue{
		client: client,
		ctx:    ctx,
	}, nil
}

// Push adds a job to the queue
func (q *RedisQueue) Push(queueName string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return q.client.RPush(q.ctx, queueName, jsonData).Err()
}

// Listen continuously listens for jobs and processes them
func (q *RedisQueue) Listen(queueName string, handler func(string) error) error {
	for {
		// Block and wait for job (timeout after 5 seconds)
		result, err := q.client.BLPop(q.ctx, 5*time.Second, queueName).Result()
		if err == redis.Nil {
			// No jobs, continue waiting
			continue
		} else if err != nil {
			return fmt.Errorf("failed to pop from queue: %w", err)
		}

		// Process the job
		jobData := result[1] // result[0] is the queue name
		if err := handler(jobData); err != nil {
			// TODO: Add to dead letter queue or retry logic
			fmt.Printf("Error processing job: %v\n", err)
		}
	}
}

// Close closes the Redis connection
func (q *RedisQueue) Close() error {
	return q.client.Close()
}
