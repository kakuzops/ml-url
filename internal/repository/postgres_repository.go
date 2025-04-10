package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/kakuzops/ml-url/internal/domain"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db      *gorm.DB
	baseURL string
}

func NewPostgresRepository(db *gorm.DB, baseURL string) *PostgresRepository {
	return &PostgresRepository{
		db:      db,
		baseURL: baseURL,
	}
}

func (r *PostgresRepository) Save(url *domain.URL) error {

	var existingURL domain.URL
	result := r.db.WithContext(context.Background()).
		Where("short_url = ?", url.ShortURL).
		First(&existingURL)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {

			return r.db.WithContext(context.Background()).Create(url).Error
		}
		return fmt.Errorf("failed to check existing URL: %w", result.Error)
	}

	existingURL.LongURL = url.LongURL
	existingURL.ExpiresAt = url.ExpiresAt

	return r.db.WithContext(context.Background()).Save(&existingURL).Error
}

func (r *PostgresRepository) FindByShortURL(shortCode string) (*domain.URL, error) {
	var url domain.URL
	result := r.db.WithContext(context.Background()).
		Where("short_url = ? AND expires_at > ?", shortCode, time.Now()).
		First(&url)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("URL not found or expired")
		}
		return nil, fmt.Errorf("failed to get URL: %w", result.Error)
	}

	return &url, nil
}

func (r *PostgresRepository) Delete(shortCode string) error {
	result := r.db.WithContext(context.Background()).
		Where("short_url = ?", shortCode).
		Delete(&domain.URL{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete URL: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("URL not found")
	}

	return nil
}

// REMOVE some urls expireds
func (r *PostgresRepository) CleanupExpiredURLs() error {
	result := r.db.WithContext(context.Background()).
		Where("expires_at < ?", time.Now()).
		Delete(&domain.URL{})

	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired URLs: %w", result.Error)
	}

	return nil
}
