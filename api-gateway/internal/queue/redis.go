package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
	ctx    context.Context
}

type CrawlJob struct {
	RunID     int    `json:"run_id"`
	ProjectID int    `json:"project_id"`
	StartURL  string `json:"start_url"`
	MaxPages  int    `json:"max_pages"`
}

func NewRedisQueue(redisURL string) (*RedisQueue, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisQueue{
		client: client,
		ctx:    ctx,
	}, nil
}

func (q *RedisQueue) PushCrawlJob(job CrawlJob) error {
	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.client.LPush(q.ctx, "crawl_queue", jobData).Err()
}

func (q *RedisQueue) Close() error {
	return q.client.Close()
}

func (q *RedisQueue) GetQueueLength(queueName string) (int64, error) {
	return q.client.LLen(q.ctx, queueName).Result()
}

func (q *RedisQueue) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	return q.client.Set(q.ctx, key, value, expiration).Err()
}

func (q *RedisQueue) Get(key string) (string, error) {
	return q.client.Get(q.ctx, key).Result()
}
