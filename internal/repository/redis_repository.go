package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kakuzops/ml-url/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client  *redis.Client
	baseURL string
}

func NewRedisRepository(client *redis.Client, baseURL string) *RedisRepository {
	return &RedisRepository{
		client:  client,
		baseURL: baseURL,
	}
}

func (r *RedisRepository) Save(url *domain.URL) error {
	data, err := json.Marshal(url)
	if err != nil {
		return fmt.Errorf("failed to marshal URL: %w", err)
	}

	ctx := context.Background()
	shortCode := strings.TrimPrefix(url.ShortURL, r.baseURL+"/")
	key := fmt.Sprintf("url:%s", shortCode)

	expiration := time.Until(url.ExpiresAt)
	if expiration < 0 {
		expiration = 0
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisRepository) FindByShortURL(shortCode string) (*domain.URL, error) {
	ctx := context.Background()
	key := fmt.Sprintf("url:%s", shortCode)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("URL not found")
		}
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	var url domain.URL
	if err := json.Unmarshal(data, &url); err != nil {
		return nil, fmt.Errorf("failed to unmarshal URL: %w", err)
	}

	return &url, nil
}

func (r *RedisRepository) Delete(shortCode string) error {
	ctx := context.Background()
	key := fmt.Sprintf("url:%s", shortCode)
	return r.client.Del(ctx, key).Err()
}
