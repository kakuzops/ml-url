package service

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
	"github.com/kakuzops/ml-url/internal/domain"
)

type URLService struct {
	repository domain.URLRepository
}

func NewURLService(repository domain.URLRepository) *URLService {
	return &URLService{
		repository: repository,
	}
}

func (s *URLService) ShortenURL(longURL string) (*domain.URL, error) {
	hash := sha256.Sum256([]byte(longURL))
	shortURL := base64.URLEncoding.EncodeToString(hash[:])[:8]

	url := &domain.URL{
		ID:        fmt.Sprintf("url_%d", time.Now().UnixNano()),
		LongURL:   longURL,
		ShortURL:  shortURL,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := s.repository.Save(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar URL: %v", err)
	}

	return url, nil
}

func (s *URLService) GetLongURL(shortURL string) (string, error) {
	url, err := s.repository.FindByShortURL(shortURL)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar URL: %v", err)
	}

	if url == nil {
		return "", fmt.Errorf("URL não encontrada")
	}

	if time.Now().After(url.ExpiresAt) {
		s.repository.Delete(shortURL)
		return "", fmt.Errorf("URL expirada")
	}

	return url.LongURL, nil
} 