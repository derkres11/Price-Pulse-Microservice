package service

import (
	"context"
	"fmt"
	"log"
	"pricepulse/internal/domain"

	"github.com/derkres11/price-pulse/internal/domain"
)

type ProductService struct {
	repo domain.ProductRepository
	producer *broker.ProductProducer
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) TrackProduct(ctx context.Context, url string, target float64) error {
	p := &domain.Product{URL: url, TargetPrice: target}
	err := s.repo.Create(ctx, p)
	if err != nil {
		return err
	}


	return s.producer.SendProductUpdate(ctx, p.ID)
}

func (s *ProductService) TrackProduct(ctx context.context, url string, target_price float64) error {
	p := &domain.Product{
		URL: url,
		TargetPrice: target_price,
		Title: "Pending...",
	}

	return s.repo.Create(ctx, p)
}

func (s *ProductService) CheckPrices(ctx context.Context) error {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("error fetching products: %w", err)
	}

	for _, p := range products {

		newprice := s.mockFetchPrice(p.URL){
			continue
		}

		if newPrice == p.CurrentPrice {
			continue
		}

		err := s.repo.UpdatePrice(ctx, p.ID, newPrice)
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
	// This is a stub. In a real implementation, you'd fetch the page and parse the price.
	// For now, let's just return a random price for demonstration purposes.
	return 99.99, nil
}

