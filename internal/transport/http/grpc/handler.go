package grpc

import (
	"context"

	"github.com/derkres11/price-pulse/internal/domain"
	desc "github.com/derkres11/price-pulse/pkg/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	desc.UnimplementedProductServiceServer
	service domain.ProductService
}

func NewHandler(svc domain.ProductService) *Handler {
	return &Handler{
		service: svc,
	}
}

func (h *Handler) GetProduct(ctx context.Context, req *desc.GetProductRequest) (*desc.GetProductResponse, error) {
	product, err := h.service.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetProductResponse{
		Id:           product.ID,
		Title:        product.Title,
		CurrentPrice: product.CurrentPrice,
		TargetPrice:  product.TargetPrice,
		CreatedAt:    timestamppb.New(product.CreatedAt),
	}, nil
}
