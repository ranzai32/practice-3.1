package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"practice4/practice-4/internal/handler"
	"practice4/practice-4/internal/repository"
	"practice4/practice-4/internal/repository/_postgres"
	"practice4/practice-4/internal/router"
	"practice4/practice-4/internal/usecase"
	"practice4/practice-4/pkg/modules"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func Run() {
	_ = godotenv.Load()

	dbConfig := initPostgreConfig()

	db := _postgres.NewPGXDialect(context.Background(), dbConfig)

	repos := repository.NewRepositories(db)
	uc := usecase.NewUserUsecase(repos.Users)
	h := handler.NewUserHandler(uc)

	apiKey := mustEnv("API_KEY")
	r := router.NewRouter(h, apiKey)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Printf("server started on :8080")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("db close error: %v", err)
	}

	log.Printf("server stopped")
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}
	return val
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func initPostgreConfig() *modules.PostgreConfig {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	return &modules.PostgreConfig{
		Host:        mustEnv("DB_HOST"),
		Port:        port,
		Username:    mustEnv("DB_USER"),
		Password:    getEnv("DB_PASSWORD", ""),
		DBName:      mustEnv("DB_NAME"),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
		ExecTimeout: 5 * time.Second,
	}
}
