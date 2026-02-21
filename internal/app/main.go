package app

import (
	"context"
	"fmt"
	"os"
	"practice4/practice-4/internal/repository"
	"practice4/practice-4/internal/repository/_postgres"
	"practice4/practice-4/pkg/modules"
	"strconv"
	"time"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()

	_postgre := _postgres.NewPGXDialect(ctx, dbConfig)

	repositories := repository.NewRepositories(_postgre)

	users, err := repositories.GetUsers()
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}
	fmt.Printf("Users: %+v\n", users)
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
		Host:        getEnv("DB_HOST", "localhost"),
		Port:        port,
		Username:    getEnv("DB_USER", "postgres"),
		Password:    getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "mydb"),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
		ExecTimeout: 5 * time.Second,
	}
}
