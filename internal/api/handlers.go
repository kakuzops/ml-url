package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	urlService URLServiceInterface
}

func NewURLHandler(urlService URLServiceInterface) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type GetURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL inválida"})
		return
	}

	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		req.URL = "http://" + req.URL
	}

	url, err := h.urlService.ShortenURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ShortenResponse{
		ShortURL: url.ShortURL,
	})
}

func (h *URLHandler) GetURLInfo(c *gin.Context) {
	shortCode := c.Param("shortURL")

	shortCode = strings.TrimPrefix(shortCode, "http://")
	shortCode = strings.TrimPrefix(shortCode, "https://")
	shortCode = strings.TrimPrefix(shortCode, "url.li/")

	urlInfo, err := h.urlService.GetURLInfo(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetURLResponse{
		ShortURL:    urlInfo.ShortURL,
		OriginalURL: urlInfo.LongURL,
		ExpiresAt:   urlInfo.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *URLHandler) RedirectToLongURL(c *gin.Context) {
	shortCode := c.Param("shortURL")

	shortCode = strings.TrimPrefix(shortCode, "http://")
	shortCode = strings.TrimPrefix(shortCode, "https://")
	shortCode = strings.TrimPrefix(shortCode, "url.li/")

	longURL, err := h.urlService.GetLongURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "http://" + longURL
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
}

func (h *URLHandler) DeleteURL(c *gin.Context) {
	shortCode := c.Param("shortURL")

	shortCode = strings.TrimPrefix(shortCode, "http://")
	shortCode = strings.TrimPrefix(shortCode, "https://")
	shortCode = strings.TrimPrefix(shortCode, "url.li/")

	if err := h.urlService.DeleteURL(shortCode); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}
