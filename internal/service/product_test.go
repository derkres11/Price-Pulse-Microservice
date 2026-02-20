package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/derkres11/price-pulse/internal/domain"
)

// repoMock matches domain.ProductRepository interface
type repoMock struct {
	products map[int64]*domain.Product
}

func (m *repoMock) Create(ctx context.Context, p *domain.Product) error {
	m.products[p.ID] = p
	return nil
}

func (m *repoMock) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (m *repoMock) UpdatePrice(ctx context.Context, id int64, newPrice float64) error {
	p, ok := m.products[id]
	if !ok {
		return errors.New("not found")
	}
	p.CurrentPrice = newPrice
	return nil
}

func (m *repoMock) GetAll(ctx context.Context) ([]*domain.Product, error) {
	var list []*domain.Product
	for _, p := range m.products {
		list = append(list, p)
	}
	return list, nil
}

// kafkaMock must match the Producer interface used in your service
type kafkaMock struct {
	sent bool
}

func (m *kafkaMock) SendProductUpdate(ctx context.Context, id int64) error {
	m.sent = true
	return nil
}

// --- TESTS ---

func TestProductService_GetByID(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := &repoMock{products: make(map[int64]*domain.Product)}

	// Seed data with your domain fields
	mockRepo.products[1] = &domain.Product{
		ID:           1,
		Title:        "Test Product",
		CurrentPrice: 100.0,
	}

	svc := NewProductService(mockRepo, nil, nil, logger)

	tests := []struct {
		name      string
		productID int64
		wantErr   bool
	}{
		{"Success", 1, false},
		{"Not Found", 999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetByID(context.Background(), tt.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestProductService_Create(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := &repoMock{products: make(map[int64]*domain.Product)}
	mockKafka := &kafkaMock{}

	svc := NewProductService(mockRepo, mockKafka, nil, logger)

	t.Run("create and notify", func(t *testing.T) {
		p := &domain.Product{ID: 10, Title: "Gadget"}
		err := svc.Create(context.Background(), p)

		if err != nil {
			t.Errorf("create failed: %v", err)
		}
		if !mockKafka.sent {
			t.Error("kafka was not notified")
		}
	})
}
