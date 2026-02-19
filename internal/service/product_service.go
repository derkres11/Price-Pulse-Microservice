package service

import (
	"context"
	"fmt"
	"log"
	"pricepulse/internal/broker"
	"pricepulse/internal/domain"
)

type ProductService struct {
	repo     domain.ProductRepository
	producer *broker.ProductProducer
}

func NewProductService(repo domain.ProductRepository, producer *broker.ProductProducer) *ProductService {
	return &ProductService{
		repo:     repo,
		producer: producer,
	}
}

func (s *ProductService) TrackProduct(ctx context.Context, url string, target_price float64) error {
	p := &domain.Product{
		URL:         url,
		TargetPrice: target_price,
		Title:       "Pending...",
	}

	err := s.repo.Create(ctx, p)
	if err != nil {
		return err
	}

	// Отправка в Kafka
	return s.producer.SendProductUpdate(ctx, p.ID)
}

func (s *ProductService) CheckPrices(ctx context.Context) error {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("error fetching products: %w", err)
	}

	for _, p := range products {
		newPrice, err := s.mockFetchPrice(p.URL)
		if err != nil {
			continue
		}

		if newPrice == p.CurrentPrice {
			continue
		}

		err = s.repo.UpdatePrice(ctx, p.ID, newPrice)
		if err != nil {
			log.Printf("error updating price for product %d: %v", p.ID, err)
			continue
		}

		if newPrice <= p.TargetPrice {
			log.Printf("Price alert for product %d! Current price: %.2f, Target price: %.2f", p.ID, newPrice, p.TargetPrice)
		}
	}
	return nil
}

func (s *ProductService) mockFetchPrice(url string) (float64, error) {
	return 99.99, nil
}
