package client

import (
	"context"
	"fmt"
	"route256/cart/internal/config"
	"route256/cart/internal/model"
	"route256/cart/internal/service"
	"route256/cart/pkg/api/product/v1"
	"route256/cart/pkg/http/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Product struct {
	product product.ProductServiceClient
}

// GetInfo implements service.ProductStorage.
func (p Product) GetInfo(ctx context.Context, sku model.SKU) (model.GetProductResponse, error) {
	token := middleware.GetAuthTokenFromContext(ctx)
	if token == "" {
		return model.GetProductResponse{}, model.ErrUnauthorized
	}
	cresp, err := p.product.GetProduct(ctx, &product.GetProductRequest{
		Token: token,
		Sku:   uint32(sku),
	})
	if err != nil {
		return model.GetProductResponse{}, mapError(ctx, err)
	}
	return model.GetProductResponse{
		Name:  cresp.Name,
		Price: cresp.Price,
	}, err
}

func NewProduct(cfg *config.GRPCClient) (*Product, error) {
	conn, err := grpc.NewClient(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("can't create grpc client connection: %w", err)
	}
	product := product.NewProductServiceClient(conn)
	return &Product{
		product: product,
	}, nil
}

var _ service.ProductStorage = Product{}
