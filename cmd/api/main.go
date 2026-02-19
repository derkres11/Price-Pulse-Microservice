package main

import (
	"context"
	"log"
	"net/http" // Standard library for HTTP server
	"os"
	"os/signal"
	"syscall"

	"github.com/derkres11/price-pulse/internal/broker"
	"github.com/derkres11/price-pulse/internal/database"
	"github.com/derkres11/price-pulse/internal/service"
	transportHTTP "github.com/derkres11/price-pulse/internal/transport/http" // Alias to avoid conflict with net/http
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Resource initialization
	dbPool := database.NewPostgresPool()
	cache := database.NewCache("localhost:6379")
	brokers := []string{"localhost:9092"}

	producer := broker.NewProductProducer(brokers, "product_updates")
	repo := database.NewProductRepo(dbPool)
	productService := service.NewProductService(repo, producer, cache)

	// Start Background Consumer (Watcher)
	consumer := broker.NewProductConsumer(brokers, "product_updates", "watcher-group")
	go func() {
		log.Println("Watcher: background consumer started")
		consumer.Start(context.Background(), func(id int64) error {
			return productService.ProcessSingleProduct(context.Background(), id)
		})
	}()

	// Initialize Handler and wrap Gin into standard http.Server
	handler := transportHTTP.NewHandler(productService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRoutes(),
	}

	// Start HTTP server in a goroutine so it doesn't block
	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error running server: %s", err.Error())
		}
	}()

	// --- SECTION: GRACEFUL SHUTDOWN ---

	// Create channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	// Notify channel on interrupt (Ctrl+C) or termination signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal
	<-quit
	log.Println("Shutdown signal received, shutting down gracefully...")
}
