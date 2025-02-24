package client

import (
	"context"
	"net/http"
	"route256/cart/internal/config"
	"route256/cart/internal/model"
	"route256/cart/internal/service"
)

type Order struct {
	client
	cfg *config.OrderClient
}

func NewOrder(cfg *config.OrderClient) Order {
	return Order{
		client: newClient(cfg.BaseURL, cfg.RequestTimeout),
		cfg:    cfg,
	}
}

type createOrderItem struct {
	SKU   model.SKU `json:"sku"`
	Count uint16    `json:"count"`
}

type createOrderRequest struct {
	User  model.UserID      `json:"user"`
	Items []createOrderItem `json:"items"`
}

type createOrderResponse struct {
	OrderID model.OrderID `json:"orderID"`
}

func (c Order) CreateOrder(ctx context.Context, req model.OrderCreateRequest) (resp model.OrderID, _ error) {
	log := getLogger(ctx, "Order.CreateOrder")

	creq := createOrderRequest{
		User:  req.UserID,
		Items: make([]createOrderItem, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		creq.Items = append(creq.Items, createOrderItem{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}

	var cresp createOrderResponse
	status, err := c.doRequest(ctx, c.cfg.CreateOrderEndpoint, &creq, &cresp)
	if err != nil {
		log.Error("can't do request", "error", err)
		return resp, model.ErrInternalError
	}

	switch status {
	case http.StatusOK, http.StatusCreated:
		return cresp.OrderID, nil
	case http.StatusPreconditionFailed:
		return resp, model.ErrPreconditionFailed
	}

	log.Error("unknown response status", "status", status)
	return resp, model.ErrInternalError
}

type getStockInfoRequest struct {
	SKU model.SKU `json:"sku"`
}

type getStockInfoResponse struct {
	Count uint64 `json:"count"`
}

func (c Order) GetStockInfo(ctx context.Context, sku model.SKU) (count uint64, _ error) {
	log := getLogger(ctx, "Order.GetStockInfo")

	creq := getStockInfoRequest{
		SKU: sku,
	}

	var cresp getStockInfoResponse
	status, err := c.doRequest(ctx, c.cfg.GetStockInfoEndpoint, &creq, &cresp)
	if err != nil {
		log.Error("can't do request", "error", err)
		return count, model.ErrInternalError
	}

	switch status {
	case http.StatusOK:
		return cresp.Count, nil
	case http.StatusNotFound:
		return count, model.ErrNotFound
	}

	log.Error("unknown response status", "status", status)
	return count, model.ErrInternalError
}

var _ service.OrderStorage = &Order{}
