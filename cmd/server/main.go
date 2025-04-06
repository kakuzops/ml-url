package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kakuzops/ml-url/internal/api"
	"github.com/kakuzops/ml-url/internal/repository"
	"github.com/kakuzops/ml-url/internal/service"
)

func main() {
	// Configuração do Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // sem senha por padrão
		DB:       0,  // usar banco padrão
	})

	// Inicialização dos componentes
	urlRepository := repository.NewRedisRepository(redisClient)
	urlService := service.NewURLService(urlRepository)
	urlHandler := api.NewURLHandler(urlService)

	// Configuração do Gin
	router := gin.Default()

	// Configuração das rotas
	router.POST("/shorten", urlHandler.ShortenURL)
	router.GET("/:shortURL", urlHandler.RedirectToLongURL)

	// Iniciar o servidor
	log.Println("Servidor iniciado na porta 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
} 