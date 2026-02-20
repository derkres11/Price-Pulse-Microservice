package domain

import (
	"context"
	"time"
)

// Product represents the core business entity of our system
// We use float64 for simplicity, but in real fintech, you'd use decimal strings or integers (cents)
type Product struct {
	ID           int64     `json:"id"`
	URL          string    `json:"url"`
	Title        string    `json:"title"`
	CurrentPrice float64   `json:"current_price"`
	TargetPrice  float64   `json:"target_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ProductRepository defines the behavior for storing and retrieving products.
// This is an Interface. It says WHAT needs to be done, but not HOW.
type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id int64) (*Product, error)
	UpdatePrice(ctx context.Context, id int64, newPrice float64) error
	GetAll(ctx context.Context) ([]*Product, error)
}

// ProductService defines the business logic operations
type ProductService interface {
	TrackProduct(ctx context.Context, url string, targetPrice float64) error
	CheckPrices(ctx context.Context) error
}

// TaskProducer defines the behavior for sending async tasks to Kafka
type TaskProducer interface {
	SendProductUpdate(ctx context.Context, id int64) error // changed name
}

// ProductCache defines the behavior for caching product data in Redis
type ProductCache interface {
	SetPrice(ctx context.Context, id int64, price float64) error // changed name
	Get(ctx context.Context, id int64) (*Product, error)
}
