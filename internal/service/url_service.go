package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/kakuzops/ml-url/internal/domain"
	"github.com/kakuzops/ml-url/internal/metrics"
)

type URLService struct {
	repo     domain.URLRepository
	baseURL  string
	duration time.Duration
}

func NewURLService(repo domain.URLRepository, baseURL string, duration time.Duration) *URLService {
	return &URLService{
		repo:     repo,
		baseURL:  baseURL,
		duration: duration,
	}
}

func (s *URLService) ShortenURL(ctx context.Context, longURL string) (*domain.URL, error) {
	if !hasProtocol(longURL) {
		longURL = "https://" + longURL
	}

	shortCode, err := generateShortCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	url := &domain.URL{
		ShortURL:  shortCode,
		LongURL:   longURL,
		ExpiresAt: time.Now().Add(s.duration),
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, url); err != nil {
		return nil, fmt.Errorf("failed to save URL: %w", err)
	}

	url.ShortURL = fmt.Sprintf("%s/%s", s.baseURL, shortCode)

	return url, nil
}

func (s *URLService) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.GetURLInfo(ctx, shortCode)
	if err != nil {
		return "", err
	}
	return url.LongURL, nil
}

func (s *URLService) GetURLInfo(ctx context.Context, shortCode string) (*domain.URL, error) {
	url, err := s.repo.FindByShortURL(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("URL not found: %w", err)
	}

	if time.Now().After(url.ExpiresAt) {
		return nil, fmt.Errorf("URL has expired")
	}

	if !hasProtocol(url.LongURL) {
		url.LongURL = "https://" + url.LongURL
	}

	if err := s.repo.Save(ctx, url); err != nil {
		return nil, fmt.Errorf("failed to update last access: %w", err)
	}

	url.ShortURL = fmt.Sprintf("%s/%s", s.baseURL, url.ShortURL)

	return url, nil
}

func generateShortCode() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:8], nil
}

func hasProtocol(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}

func (s *URLService) DeleteURL(ctx context.Context, shortCode string) error {

	_, err := s.repo.FindByShortURL(ctx, shortCode)
	if err != nil {
		return fmt.Errorf("URL not found: %w", err)
	}

	if err := s.repo.Delete(ctx, shortCode); err != nil {
		return fmt.Errorf("failed to delete URL: %w", err)
	}

	metrics.DecrementActiveURLs()

	return nil
}
