package handler

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/pkg/api/stock/v1"
)

type StockService interface {
	GetStockInfo(ctx context.Context, sku model.SKU) (count uint64, err error)
}

type Stock struct {
	stock.UnimplementedStockServer
	service StockService
}

func NewStock(service StockService) *Stock {
	return &Stock{service: service}
}

func (s Stock) GetInfo(ctx context.Context, req *stock.GetInfoRequest) (*stock.GetInfoResponse, error) {
	count, err := s.service.GetStockInfo(ctx, model.SKU(req.Sku))
	if err != nil {
		return nil, err
	}
	return &stock.GetInfoResponse{Count: count}, nil
}
