package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kakuzops/ml-url/internal/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CachedRepository struct {
	db       *gorm.DB
	redis    *redis.Client
	baseURL  string
	cacheTTL time.Duration
}

func NewCachedRepository(db *gorm.DB, redisClient *redis.Client, baseURL string, cacheTTL time.Duration) *CachedRepository {
	return &CachedRepository{
		db:       db,
		redis:    redisClient,
		baseURL:  baseURL,
		cacheTTL: cacheTTL,
	}
}

func (r *CachedRepository) Save(url *domain.URL) error {

	if err := r.saveToDatabase(url); err != nil {
		return err
	}

	return r.saveToCache(url)
}

func (r *CachedRepository) FindByShortURL(shortCode string) (*domain.URL, error) {

	url, err := r.findInCache(shortCode)
	if err == nil {
		return url, nil
	}

	url, err = r.findInDatabase(shortCode)
	if err != nil {
		return nil, err
	}

	if url.ExpiresAt.After(time.Now()) {
		if err := r.saveToCache(url); err != nil {

			fmt.Printf("Failed to save to cache: %v\n", err)
		}
	}

	return url, nil
}

func (r *CachedRepository) Delete(shortCode string) error {

	if err := r.deleteFromDatabase(shortCode); err != nil {
		return err
	}

	return r.deleteFromCache(shortCode)
}

func (r *CachedRepository) saveToDatabase(url *domain.URL) error {
	var existingURL domain.URL
	result := r.db.Where("short_url = ?", url.ShortURL).First(&existingURL)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return r.db.Create(url).Error
		}
		return fmt.Errorf("failed to check existing URL: %w", result.Error)
	}

	existingURL.LongURL = url.LongURL
	existingURL.ExpiresAt = url.ExpiresAt

	return r.db.Save(&existingURL).Error
}

func (r *CachedRepository) findInDatabase(shortCode string) (*domain.URL, error) {
	var url domain.URL
	result := r.db.Where("short_url = ? AND expires_at > ?", shortCode, time.Now()).First(&url)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("URL not found or expired")
		}
		return nil, fmt.Errorf("failed to get URL: %w", result.Error)
	}

	return &url, nil
}

func (r *CachedRepository) deleteFromDatabase(shortCode string) error {
	result := r.db.Where("short_url = ?", shortCode).Delete(&domain.URL{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete URL: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("URL not found")
	}
	return nil
}

func (r *CachedRepository) saveToCache(url *domain.URL) error {
	data, err := json.Marshal(url)
	if err != nil {
		return fmt.Errorf("failed to marshal URL: %w", err)
	}

	key := r.getCacheKey(url.ShortURL)
	ctx := context.Background()

	ttl := time.Until(url.ExpiresAt)
	if ttl > r.cacheTTL {
		ttl = r.cacheTTL
	}

	return r.redis.Set(ctx, key, data, ttl).Err()
}

func (r *CachedRepository) findInCache(shortCode string) (*domain.URL, error) {
	key := r.getCacheKey(shortCode)
	ctx := context.Background()

	data, err := r.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("URL not found in cache")
		}
		return nil, fmt.Errorf("failed to get URL from cache: %w", err)
	}

	var url domain.URL
	if err := json.Unmarshal(data, &url); err != nil {
		return nil, fmt.Errorf("failed to unmarshal URL: %w", err)
	}

	return &url, nil
}

func (r *CachedRepository) deleteFromCache(shortCode string) error {
	key := r.getCacheKey(shortCode)
	ctx := context.Background()
	return r.redis.Del(ctx, key).Err()
}

func (r *CachedRepository) getCacheKey(shortCode string) string {
	return fmt.Sprintf("url:%s", shortCode)
}
