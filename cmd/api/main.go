package main

import (
	"log"

	"github.com/derkres11/price-pulse/internal/transport/http"
)

func main() {
	handler := http.NewHandler()
	srv := handler.InitRoutes()

	log.Println("Server started on :8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}
