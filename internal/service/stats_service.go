package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLStats struct {
	ShortURL    string    `json:"short_url"`
	LongURL     string    `json:"long_url"`
	AccessCount int64     `json:"access_count"`
	LastAccess  time.Time `json:"last_access"`
}

type StatsService struct {
	redis *redis.Client
}

func NewStatsService(redis *redis.Client) *StatsService {
	return &StatsService{
		redis: redis,
	}
}

func (s *StatsService) IncrementAccess(shortURL, longURL string) error {
	ctx := context.Background()
	key := fmt.Sprintf("stats:url:%s", shortURL)
	pipe := s.redis.Pipeline()
	pipe.HIncrBy(ctx, key, "access_count", 1)
	pipe.HSet(ctx, key, "last_access", time.Now().Format(time.RFC3339))
	pipe.HSet(ctx, key, "long_url", longURL)
	pipe.Expire(ctx, key, 30*24*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *StatsService) GetTopURLs(limit int) ([]URLStats, error) {
	ctx := context.Background()

	pattern := "stats:url:*"
	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()

	stats := make([]URLStats, 0)

	for iter.Next(ctx) {
		key := iter.Val()

		data, err := s.redis.HGetAll(ctx, key).Result()
		if err != nil {
			continue
		}

		count, _ := s.redis.HGet(ctx, key, "access_count").Int64()

		lastAccess, _ := time.Parse(time.RFC3339, data["last_access"])

		stats = append(stats, URLStats{
			ShortURL:    key[10:],
			LongURL:     data["long_url"],
			AccessCount: count,
			LastAccess:  lastAccess,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].AccessCount > stats[j].AccessCount
	})

	if len(stats) > limit {
		stats = stats[:limit]
	}

	return stats, nil
}

func (s *StatsService) GetURLStats(shortURL string) (*URLStats, error) {
	ctx := context.Background()
	key := fmt.Sprintf("stats:url:%s", shortURL)

	data, err := s.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get URL stats: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("URL stats not found")
	}

	count, _ := s.redis.HGet(ctx, key, "access_count").Int64()
	lastAccess, _ := time.Parse(time.RFC3339, data["last_access"])

	return &URLStats{
		ShortURL:    shortURL,
		LongURL:     data["long_url"],
		AccessCount: count,
		LastAccess:  lastAccess,
	}, nil
}
