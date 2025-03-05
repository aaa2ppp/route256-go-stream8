package client

import (
	"context"
	"fmt"
	"route256/cart/internal/config"
	"route256/cart/internal/model"
	"route256/cart/internal/service"
	"route256/loms/pkg/api/order/v1"
	"route256/loms/pkg/api/stock/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Order struct {
	order order.OrderClient
	stock stock.StockClient
}

// CreateOrder implements service.OrderStorage.
func (o *Order) CreateOrder(ctx context.Context, req model.OrderCreateRequest) (model.OrderID, error) {
	creq := &order.CreateRequest{
		User:  int64(req.UserID),
		Items: make([]*order.Item, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		creq.Items = append(creq.Items, &order.Item{
			Sku:   int32(item.SKU),
			Count: uint32(item.Count),
		})
	}
	cresp, err := o.order.Create(ctx, creq)
	if err != nil {
		return 0, mapError(ctx, err)
	}
	return model.OrderID(cresp.OrderID), nil
}

// GetStockInfo implements service.OrderStorage.
func (o *Order) GetStockInfo(ctx context.Context, sku model.SKU) (count uint64, _ error) {
	cresp, err := o.stock.GetInfo(ctx, &stock.GetInfoRequest{Sku: int32(sku)})
	if err != nil {
		return 0, mapError(ctx, err)
	}
	return cresp.Count, nil
}

func NewOrder(cfg *config.GRPCLOMSClient) (*Order, error) {
	conn, err := grpc.NewClient(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("can't create grpc client connection: %w", err)
	}
	order := order.NewOrderClient(conn)
	stock := stock.NewStockClient(conn)
	return &Order{
		order: order,
		stock: stock,
	}, nil
}

var _ service.OrderStorage = &Order{}
