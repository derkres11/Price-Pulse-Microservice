package main

import (
	"log"
	"pricepulse/internal/database"

	"github.com/derkres11/price-pulse/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db := database.NewPostgresPool()
	defer db.Close()

	log.Println("Server is starting...")
	handler := http.NewHandler(db)
	srv := handler.InitRoutes()

	log.Println("Server started on :8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}
