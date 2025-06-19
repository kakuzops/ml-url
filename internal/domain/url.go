package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URL struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	LongURL   string         `json:"long_url" gorm:"type:text;not null"`
	ShortURL  string         `json:"short_url" gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null;index"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (URL) TableName() string {
	return "shorten_url"
}

func (u *URL) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type URLRepository interface {
	Save(ctx context.Context, url *URL) error
	FindByShortURL(ctx context.Context, shortURL string) (*URL, error)
	Delete(ctx context.Context, shortURL string) error
}
