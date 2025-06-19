package api

import (
	"context"

	"github.com/kakuzops/ml-url/internal/domain"
)

type URLServiceInterface interface {
	ShortenURL(ctx context.Context, longURL string) (*domain.URL, error)
	GetLongURL(ctx context.Context, shortCode string) (string, error)
	GetURLInfo(ctx context.Context, shortCode string) (*domain.URL, error)
	DeleteURL(ctx context.Context, shortCode string) error
}
