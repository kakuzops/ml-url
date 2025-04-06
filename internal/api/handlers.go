package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/kakuzops/ml-url/internal/service"
)

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
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

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL inválida"})
		return
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

func (h *URLHandler) RedirectToLongURL(c *gin.Context) {
	shortURL := c.Param("shortURL")

	longURL, err := h.urlService.GetLongURL(shortURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
} 