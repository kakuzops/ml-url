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
	cfg := config.LoadConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	urlRepo := repository.NewRedisRepository(redisClient, cfg.BaseURL)

	urlService := service.NewURLService(urlRepo, cfg.BaseURL, cfg.Duration)

	handlers := api.NewURLHandler(urlService)

	router := gin.Default()

	router.Use(metrics.MetricsMiddleware())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	router.POST("/shorten", handlers.ShortenURL)
	router.GET("/:shortURL", handlers.RedirectToLongURL)
	router.GET("/info/:shortURL", handlers.GetURLInfo)
	router.DELETE("/:shortURL", handlers.DeleteURL)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
