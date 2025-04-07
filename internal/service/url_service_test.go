package service

import (
	"strings"
	"testing"
	"time"
	"github.com/kakuzops/ml-url/internal/domain"
)

type mockRepository struct {
	urls map[string]*domain.URL
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		urls: make(map[string]*domain.URL),
	}
}

func (m *mockRepository) Save(url *domain.URL) error {
	shortCode := strings.TrimPrefix(url.ShortURL, "http://url.li/")
	m.urls[shortCode] = url
	return nil
}

func (m *mockRepository) FindByShortURL(shortCode string) (*domain.URL, error) {
	url, exists := m.urls[shortCode]
	if !exists {
		return nil, nil
	}
	return url, nil
}

func (m *mockRepository) Delete(shortCode string) error {
	delete(m.urls, shortCode)
	return nil
}

func TestShortenURL(t *testing.T) {
	repo := newMockRepository()
	service := NewURLService(repo)

	longURL := "https://www.google.com.br"
	url, err := service.ShortenURL(longURL)

	if err != nil {
		t.Errorf("Erro inesperado ao encurtar URL: %v", err)
	}

	if url.LongURL != longURL {
		t.Errorf("URL longa esperada %s, obtida %s", longURL, url.LongURL)
	}

	if !strings.HasPrefix(url.ShortURL, "http://url.li/") {
		t.Errorf("URL curta deve começar com http://url.li/, obtida %s", url.ShortURL)
	}

	shortCode := strings.TrimPrefix(url.ShortURL, "http://url.li/")
	if len(shortCode) != 8 {
		t.Errorf("Código da URL curta deve ter 8 caracteres, obtido %d", len(shortCode))
	}
}

func TestGetLongURL(t *testing.T) {
	repo := newMockRepository()
	service := NewURLService(repo)

	longURL := "https://www.google.com.br"
	url, _ := service.ShortenURL(longURL)
	shortCode := strings.TrimPrefix(url.ShortURL, "http://url.li/")

	retrievedURL, err := service.GetLongURL(shortCode)
	if err != nil {
		t.Errorf("Erro inesperado ao recuperar URL: %v", err)
	}

	if retrievedURL != longURL {
		t.Errorf("URL longa esperada %s, obtida %s", longURL, retrievedURL)
	}
}

func TestGetExpiredURL(t *testing.T) {
	repo := newMockRepository()
	service := NewURLService(repo)

	shortCode := "expired"
	url := &domain.URL{
		ID:        "test_id",
		LongURL:   "https://www.google.com.br",
		ShortURL:  "http://url.li/" + shortCode,
		CreatedAt: time.Now().Add(-25 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	repo.Save(url)

	_, err := service.GetLongURL(shortCode)
	if err == nil {
		t.Error("Esperado erro de URL expirada, mas nenhum erro foi retornado")
	}
} 