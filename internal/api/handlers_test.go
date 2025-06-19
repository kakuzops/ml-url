package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kakuzops/ml-url/internal/domain"
)

type mockURLService struct {
	urls map[string]*domain.URL
}

func newMockURLService() *mockURLService {
	return &mockURLService{
		urls: make(map[string]*domain.URL),
	}
}

func (m *mockURLService) ShortenURL(ctx context.Context, longURL string) (*domain.URL, error) {
	url := &domain.URL{
		LongURL:   longURL,
		ShortURL:  "testshort",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	m.urls[url.ShortURL] = url
	return url, nil
}

func (m *mockURLService) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	if url, exists := m.urls[shortCode]; exists {
		return url.LongURL, nil
	}
	return "", fmt.Errorf("URL not found")
}

func (m *mockURLService) GetURLInfo(ctx context.Context, shortCode string) (*domain.URL, error) {
	if url, exists := m.urls[shortCode]; exists {
		return url, nil
	}
	return nil, fmt.Errorf("URL not found")
}

func (m *mockURLService) DeleteURL(ctx context.Context, shortCode string) error {
	if _, exists := m.urls[shortCode]; !exists {
		return fmt.Errorf("URL not found")
	}
	delete(m.urls, shortCode)
	return nil
}

func TestDeleteURL(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockService := newMockURLService()
	handler := NewURLHandler(mockService, nil)
	router := gin.New()
	router.DELETE("/:shortURL", handler.DeleteURL)

	t.Run("Delete existing URL", func(t *testing.T) {
		url, _ := mockService.ShortenURL(context.Background(), "https://www.example.com")
		shortCode := url.ShortURL

		req := httptest.NewRequest("DELETE", "/"+shortCode, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Delete non-existing URL", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/naoexiste", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
