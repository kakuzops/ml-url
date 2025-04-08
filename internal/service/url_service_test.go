package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kakuzops/ml-url/internal/domain"
)

type mockRepository struct {
	urls    map[string]*domain.URL
	baseURL string
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		urls:    make(map[string]*domain.URL),
		baseURL: "http://url.li",
	}
}

func (m *mockRepository) Save(url *domain.URL) error {
	shortCode := strings.TrimPrefix(url.ShortURL, m.baseURL+"/")
	m.urls[shortCode] = url
	return nil
}

func (m *mockRepository) FindByShortURL(shortCode string) (*domain.URL, error) {
	url, exists := m.urls[shortCode]
	if !exists {
		return nil, fmt.Errorf("URL not found")
	}
	return url, nil
}

func (m *mockRepository) Delete(shortCode string) error {
	delete(m.urls, shortCode)
	return nil
}

func TestShortenURL(t *testing.T) {
	repo := newMockRepository()
	service := NewURLService(repo, "http://url.li", 24*time.Hour)

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
	service := NewURLService(repo, "http://url.li", 24*time.Hour)

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
	service := NewURLService(repo, "http://url.li", 24*time.Hour)

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

func TestDeleteURL(t *testing.T) {

	repo := newMockRepository()
	service := NewURLService(repo, "http://url.li", 24*time.Hour)

	t.Run("Delete existing URL", func(t *testing.T) {
		longURL := "https://www.example.com"
		url, err := service.ShortenURL(longURL)
		if err != nil {
			t.Fatalf("Erro inesperado ao criar URL: %v", err)
		}

		shortCode := strings.TrimPrefix(url.ShortURL, "http://url.li/")

		err = service.DeleteURL(shortCode)
		if err != nil {
			t.Errorf("Erro inesperado ao deletar URL: %v", err)
		}

		_, err = service.GetURLInfo(shortCode)
		if err == nil {
			t.Error("URL ainda existe após deleção")
		}
	})

	t.Run("Delete non-existing URL", func(t *testing.T) {
		err := service.DeleteURL("naoexiste")
		if err == nil {
			t.Error("Esperado erro ao deletar URL inexistente, mas nenhum erro foi retornado")
		}
	})

	t.Run("Delete already deleted URL", func(t *testing.T) {
		longURL := "https://www.example.com"
		url, err := service.ShortenURL(longURL)
		if err != nil {
			t.Fatalf("Erro inesperado ao criar URL: %v", err)
		}

		shortCode := strings.TrimPrefix(url.ShortURL, "http://url.li/")

		err = service.DeleteURL(shortCode)
		if err != nil {
			t.Errorf("Erro inesperado ao deletar URL: %v", err)
		}

		err = service.DeleteURL(shortCode)
		if err == nil {
			t.Error("Esperado erro ao deletar URL já deletada, mas nenhum erro foi retornado")
		}
	})
}
