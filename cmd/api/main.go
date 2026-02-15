package main

import (
	"log"
	"pricepulse/internal/database"

	"github.com/derkres11/price-pulse/internal/service"
	"github.com/derkres11/price-pulse/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// 1. Db Connection Pool
	dbPool := database.NewPostgresPool()
	defer dbPool.Close()

	// 2. Product Repository
	repo := database.NewProductRepo(dbPool)

	// 3. Product Service
	productService := service.NewProductService(repo)

	// 4. HTTP Handler
	handler := http.NewHandler(productService)

	// 5. Start HTTP Server
	srv := handler.InitRoutes()
	log.Println("Server started on :8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}
