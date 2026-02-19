package main

import (
	"context"
	"log"
	"pricepulse/internal/broker"
	"pricepulse/internal/database"
	"pricepulse/internal/service"
	"pricepulse/internal/transport/http"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbPool := database.NewPostgresPool()
	defer dbPool.Close()

	cache := database.NewCache("localhost:6379")

	brokers := []string{"localhost:9092"}

	producer := broker.NewProductProducer(brokers, "product_updates")
	defer producer.Close()

	repo := database.NewProductRepo(dbPool)

	productService := service.NewProductService(repo, producer, cache)

	consumer := broker.NewProductConsumer(brokers, "product_updates", "watcher-group")
	go func() {
		log.Println("Watcher: background consumer started")
		consumer.Start(context.Background(), func(id int64) error {
			return productService.ProcessSingleProduct(context.Background(), id)
		})
	}()

	handler := http.NewHandler(productService)
	srv := handler.InitRoutes()

	log.Println("Server started on :8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}
