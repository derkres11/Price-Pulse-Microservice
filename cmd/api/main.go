package main

import (
	"log"
	"pricepulse/internal/database"

	"github.com/derkres11/price-pulse/internal/transport/http"
)

func main() {
	db := database.NewPostgresPool()
	handler := http.NewHandler(db)
	srv := handler.InitRoutes()

	log.Println("Server started on :8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}
