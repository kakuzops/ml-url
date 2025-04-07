package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kakuzops/ml-url/internal/api"
	"github.com/kakuzops/ml-url/internal/metrics"
	"github.com/kakuzops/ml-url/internal/repository"
	"github.com/kakuzops/ml-url/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	urlRepository := repository.NewRedisRepository(redisClient)
	urlService := service.NewURLService(urlRepository)
	urlHandler := api.NewURLHandler(urlService)

	router := gin.Default()

	router.Use(metrics.MetricsMiddleware())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.POST("/shorten", urlHandler.ShortenURL)
	router.GET("/info/:shortURL", urlHandler.GetURLInfo)
	router.GET("/:shortURL", urlHandler.RedirectToLongURL)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	log.Println("Servidor iniciado na porta 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
} 