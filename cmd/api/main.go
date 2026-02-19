package main

import (
	"context"
	"log"
	"pricepulse/internal/database"

	"github.com/derkres11/price-pulse/internal/service"
	"github.com/derkres11/price-pulse/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// 1. Инфраструктура
	dbPool := database.NewPostgresPool()
	defer dbPool.Close()

	// 2. Брокер (Producer)
	producer := broker.NewProductProducer([]string{"localhost:9092"}, "product_updates")
	defer producer.Close()

	// 3. Слои приложения
	repo := database.NewProductRepo(dbPool)
	svc := service.NewProductService(repo, producer)

	// 4. ЗАПУСК КОНСЬЮМЕРА (Watcher)
	// Запускаем в фоне через 'go func()'
	consumer := broker.NewProductConsumer([]string{"localhost:9092"}, "product_updates", "watcher-group")
	go func() {
		log.Println("Watcher: Kafka Consumer is running...")
		consumer.Start(context.Background(), func(id int64) error {
			// Здесь вызываем метод сервиса, который мы подготовили
			return svc.ProcessSingleProduct(context.Background(), id)
		})
	}()

	// 5. Запуск API
	handler := http.NewHandler(svc)
	log.Println("API: Server is running on :8080")
	handler.InitRoutes().Run(":8080")
}
