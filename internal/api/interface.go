package api

import (
	"github.com/kakuzops/ml-url/internal/domain"
)

type URLServiceInterface interface {
	ShortenURL(longURL string) (*domain.URL, error)
	GetLongURL(shortCode string) (string, error)
	GetURLInfo(shortCode string) (*domain.URL, error)
	DeleteURL(shortCode string) error
}
