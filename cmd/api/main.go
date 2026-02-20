package main

import (
	"context"
	"log/slog" // New structured logging package
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net"

	"github.com/derkres11/price-pulse/internal/broker"
	"github.com/derkres11/price-pulse/internal/database"
	"github.com/derkres11/price-pulse/internal/service"
	grpcHandler "github.com/derkres11/price-pulse/internal/transport/grpc"
	transportHTTP "github.com/derkres11/price-pulse/internal/transport/http"
	desc "github.com/derkres11/price-pulse/pkg/api/v1"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Initialize structured logger (JSON format)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger) // Set as global logger

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment variables")
	}

	// Resource initialization
	dbPool := database.NewPostgresPool()
	cache := database.NewCache("localhost:6379")
	brokers := []string{"localhost:9092"}

	producer := broker.NewProductProducer(brokers, "product_updates")
	repo := database.NewProductRepo(dbPool)
	productService := service.NewProductService(repo, producer, cache, logger)

	// Start Background Consumer (Watcher)
	consumer := broker.NewProductConsumer(brokers, "product_updates", "watcher-group")
	go func() {
		slog.Info("Watcher: background consumer started")
		consumer.Start(context.Background(), func(id int64) error {
			return productService.ProcessSingleProduct(context.Background(), id)
		})
	}()

	// Initialize Handler and wrap Gin into standard http.Server
	handler := transportHTTP.NewHandler(productService, logger)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRoutes(),
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("failed to listen for gRPC", "error", err)
		os.Exit(1)
	}

	sServer := grpc.NewServer()
	desc.RegisterProductServiceServer(sServer, grpcHandler.NewHandler(productService))

	go func() {
		slog.Info("gRPC server started", slog.String("port", "50051"))
		if err := sServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed", "error", err)
		}
	}()

	// Start HTTP server in a goroutine
	go func() {
		slog.Info("Server started", slog.String("port", "8080"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to run server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// --- SECTION: GRACEFUL SHUTDOWN ---

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Shutdown HTTP server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	// 2. Close Kafka Producer
	if err := producer.Close(); err != nil {
		slog.Error("Kafka producer close error", slog.String("error", err.Error()))
	}

	// 3. Close Kafka Consumer
	if err := consumer.Close(); err != nil {
		slog.Error("Kafka consumer close error", slog.String("error", err.Error()))
	}

	// 4. Close Database connection pool
	dbPool.Close()

	slog.Info("Server exited properly")
}
