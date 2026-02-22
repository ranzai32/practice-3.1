package main

import "practice4/practice-4/internal/app"

// @title Practice4 API
// @version 1.0
// @description REST API with layered architecture
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY
func main() {
	app.Run()
}
