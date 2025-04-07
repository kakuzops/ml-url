package repository

import (
	"context"
	"encoding/json"
	"strings"
	"github.com/go-redis/redis/v8"
	"github.com/kakuzops/ml-url/internal/domain"
)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisRepository) Save(url *domain.URL) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	shortCode := strings.TrimPrefix(url.ShortURL, "http://url.li/")

	err = r.client.Set(r.ctx, shortCode, data, url.ExpiresAt.Sub(url.CreatedAt)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) FindByShortURL(shortCode string) (*domain.URL, error) {
	data, err := r.client.Get(r.ctx, shortCode).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var url domain.URL
	err = json.Unmarshal(data, &url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *RedisRepository) Delete(shortCode string) error {
	return r.client.Del(r.ctx, shortCode).Err()
} 