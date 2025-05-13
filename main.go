package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"posgoexpert-rate-limiter/internal/infra/database"
	"posgoexpert-rate-limiter/internal/limiter"
	"posgoexpert-rate-limiter/internal/logger"
	"posgoexpert-rate-limiter/internal/middleware"

	"github.com/joho/godotenv"
)

func main() {
	// Carrega variáveis de ambiente do .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, seguindo com variáveis de ambiente do sistema.")
	}

	logger.Init()
	// Conecta ao Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0 // default DB 0

	store := database.NewRedisStore(&database.RedisConfig{
		Addr:     redisAddr,
		Password: redisPassword,
		Db:       redisDB,
	})

	// Limites padrão
	defaultLimit := 10
	defaultExpiration := 60 // segundos

	limiterService := limiter.NewLimiter(store, defaultLimit, defaultExpiration)

	// Middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	handler := middleware.RateLimitMiddleware(limiterService)(mux)

	fmt.Println("Servidor iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
