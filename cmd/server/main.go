package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"github.com/kakuzops/ml-url/internal/api"
	"github.com/kakuzops/ml-url/internal/config"
	"github.com/kakuzops/ml-url/internal/metrics"
	"github.com/kakuzops/ml-url/internal/repository"
	"github.com/kakuzops/ml-url/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repository
	urlRepo := repository.NewRedisRepository(redisClient, cfg.BaseURL)

	// Initialize service
	urlService := service.NewURLService(urlRepo, cfg.BaseURL, cfg.Duration)

	// Initialize handlers
	handlers := api.NewURLHandler(urlService)

	// Initialize router
	router := gin.Default()

	// Add metrics middleware
	router.Use(metrics.MetricsMiddleware())

	// Add Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// Add routes
	router.POST("/shorten", handlers.ShortenURL)
	router.GET("/:shortURL", handlers.RedirectToLongURL)
	router.GET("/info/:shortURL", handlers.GetURLInfo)

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
