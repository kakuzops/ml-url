package service

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"github.com/kakuzops/ml-url/internal/domain"
)

const baseURL = "http://url.li/"

type URLService struct {
	repository domain.URLRepository
}

func NewURLService(repository domain.URLRepository) *URLService {
	return &URLService{
		repository: repository,
	}
}

func (s *URLService) ShortenURL(longURL string) (*domain.URL, error) {
	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "http://" + longURL
	}

	hash := sha256.Sum256([]byte(longURL))
	shortCode := base64.URLEncoding.EncodeToString(hash[:])[:8]
	shortURL := baseURL + shortCode

	url := &domain.URL{
		ID:        fmt.Sprintf("url_%d", time.Now().UnixNano()),
		LongURL:   longURL,
		ShortURL:  shortURL,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	existingURL, err := s.repository.FindByShortURL(shortCode)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar URL existente: %v", err)
	}

	if existingURL != nil && time.Now().Before(existingURL.ExpiresAt) {
		return existingURL, nil
	}

	err = s.repository.Save(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar URL: %v", err)
	}

	return url, nil
}

func (s *URLService) GetURLInfo(shortCode string) (*domain.URL, error) {
	url, err := s.repository.FindByShortURL(shortCode)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar URL: %v", err)
	}

	if url == nil {
		return nil, fmt.Errorf("URL não encontrada")
	}

	if time.Now().After(url.ExpiresAt) {
		s.repository.Delete(shortCode)
		return nil, fmt.Errorf("URL expirada")
	}

	// Garantir que a URL tenha o protocolo
	if !strings.HasPrefix(url.LongURL, "http://") && !strings.HasPrefix(url.LongURL, "https://") {
		url.LongURL = "http://" + url.LongURL
	}

	return url, nil
}

func (s *URLService) GetLongURL(shortCode string) (string, error) {
	url, err := s.GetURLInfo(shortCode)
	if err != nil {
		return "", err
	}

	return url.LongURL, nil
} 