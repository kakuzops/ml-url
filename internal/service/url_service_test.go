package service

import (
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
	m.urls[url.ShortURL] = url
	return nil
}

func (m *mockRepository) FindByShortURL(shortURL string) (*domain.URL, error) {
	url, exists := m.urls[shortURL]
	if !exists {
		return nil, nil
	}
	return url, nil
}

func (m *mockRepository) Delete(shortURL string) error {
	delete(m.urls, shortURL)
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

	if len(url.ShortURL) != 8 {
		t.Errorf("Tamanho da URL curta esperado 8, obtido %d", len(url.ShortURL))
	}
}

func TestGetLongURL(t *testing.T) {
	repo := newMockRepository()
	service := NewURLService(repo)

	// Criar uma URL curta
	longURL := "https://www.google.com.br"
	url, _ := service.ShortenURL(longURL)

	retrievedURL, err := service.GetLongURL(url.ShortURL)
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

	// Criar uma URL expirada
	url := &domain.URL{
		ID:        "test_id",
		LongURL:   "https://www.google.com.br",
		ShortURL:  "expired",
		CreatedAt: time.Now().Add(-25 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	repo.Save(url)

	// Tentar recuperar a URL expirada
	_, err := service.GetLongURL(url.ShortURL)
	if err == nil {
		t.Error("Esperado erro de URL expirada, mas nenhum erro foi retornado")
	}
} 