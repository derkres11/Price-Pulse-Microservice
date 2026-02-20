package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/derkres11/price-pulse/internal/domain"
)

// repoMock simulates the database repository
type repoMock struct {
	products map[int64]*domain.Product
}

// Create mimics saving a product
func (m *repoMock) Create(ctx context.Context, p *domain.Product) error {
	m.products[p.ID] = p
	return nil
}

// GetByID mimics retrieving a product
func (m *repoMock) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return p, nil
}

// UpdatePrice mimics updating the price
func (m *repoMock) UpdatePrice(ctx context.Context, id int64, price float64) error {
	if p, ok := m.products[id]; ok {
		// Assuming the field is actually 'Price'.
		// If it's different in your domain, change it here.
		p.CurrentPrice = price
		return nil
	}
	return errors.New("not found")
}

// GetAll implements the missing method required by domain.ProductRepository
func (m *repoMock) GetAll(ctx context.Context) ([]*domain.Product, error) {
	var list []*domain.Product
	for _, p := range m.products {
		list = append(list, p)
	}
	return list, nil
}

func TestProductService_GetByID(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mock := &repoMock{
		products: make(map[int64]*domain.Product),
	}

	testID := int64(42)
	// Make sure fields (ID, Title, Price) match your domain.Product struct exactly
	testProduct := &domain.Product{
		ID:           testID,
		Title:        "Smartphone",
		CurrentPrice: 599.99,
	}

	mock.products[testID] = testProduct
	svc := NewProductService(mock, nil, nil, logger)

	t.Run("successful retrieval", func(t *testing.T) {
		res, err := svc.GetByID(context.Background(), testID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if res.Title != "Smartphone" {
			t.Errorf("Expected title 'Smartphone', got %s", res.Title)
		}
	})

	t.Run("product not found", func(t *testing.T) {
		_, err := svc.GetByID(context.Background(), 999)
		if err == nil {
			t.Error("Expected error for non-existent product, got nil")
		}
	})
}
