package domain

import (
	"time"
)

type URL struct {
	ID        string    `json:"id"`
	LongURL   string    `json:"long_url"`
	ShortURL  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type URLRepository interface {
	Save(url *URL) error
	FindByShortURL(shortURL string) (*URL, error)
	Delete(shortURL string) error
}
