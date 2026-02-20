package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/derkres11/price-pulse/internal/domain"
)

type ProductService struct {
	repo     domain.ProductRepository
	producer domain.TaskProducer
	cache    domain.ProductCache
	logger   *slog.Logger
}

func NewProductService(
	repo domain.ProductRepository,
	producer domain.TaskProducer,
	cache domain.ProductCache,
	logger *slog.Logger,
) *ProductService {
	return &ProductService{
		repo:     repo,
		producer: producer,
		cache:    cache,
		logger:   logger,
	}
}

func (s *ProductService) Create(ctx context.Context, p *domain.Product) error {
	s.logger.Info("creating new product", slog.String("url", p.URL))

	//Save to DB
	if err := s.repo.Create(ctx, p); err != nil {
		s.logger.Error("failed to create product in db",
			slog.String("error", err.Error()),
			slog.String("url", p.URL))
		return err
	}

	//Sending to Kafka
	if err := s.producer.SendProductUpdate(ctx, p.ID); err != nil {
		s.logger.Error("failed to send kafka notification",
			slog.Int64("id", p.ID),
			slog.String("error", err.Error()))

	}

	return nil
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

// ProcessSingleProduct is the core logic for the Watcher
func (s *ProductService) ProcessSingleProduct(ctx context.Context, id int64) error {
	log.Printf("Watcher: processing product %d", id)
	newPrice, _ := s.mockFetchPrice("url")

	_ = s.cache.SetPrice(ctx, id, newPrice)
	return s.repo.UpdatePrice(ctx, id, newPrice)
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	s.logger.Info("fetching product", slog.Int64("id", id))

	// Trying to get from Cache
	cachedProduct, err := s.cache.Get(ctx, id)
	if err == nil && cachedProduct != nil {
		s.logger.Debug("cache hit", slog.Int64("id", id))
		return cachedProduct, nil
	}

	// Going to DB if not found in cache
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to find product in db", slog.Int64("id", id), slog.String("error", err.Error()))
		return nil, err
	}

	// Caching
	s.logger.Debug("cache miss, loading from db", slog.Int64("id", id))

	return product, nil
}
