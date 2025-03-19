package handler

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/pkg/api/order/v1"
)

type LomsService interface {
	CreateOrder(ctx context.Context, req model.CreateOrderRequest) (orderID model.OrderID, err error)
	GetOrderInfo(ctx context.Context, orderID model.OrderID) (resp model.Order, err error)
	PayOrder(ctx context.Context, orderID model.OrderID) error
	CancelOrder(ctx context.Context, orderID model.OrderID) error
	GetStockInfo(ctx context.Context, sku model.SKU) (count uint64, err error)
}

type Order struct {
	order.UnimplementedOrderServer
	service LomsService
}

func NewOrder(service LomsService) *Order {
	return &Order{service: service}
}

func (o Order) Create(ctx context.Context, req *order.CreateRequest) (*order.CreateResponse, error) {
	mreq := model.CreateOrderRequest{
		UserID: model.UserID(req.User),
		Items:  make([]model.OrderItem, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		mreq.Items = append(mreq.Items, model.OrderItem{
			SKU:   model.SKU(item.Sku),
			Count: int(item.Count),
		})
	}
	orderID, err := o.service.CreateOrder(ctx, mreq)
	if err != nil {
		return nil, mapError(ctx, err)
	}
	return &order.CreateResponse{OrderID: int64(orderID)}, nil
}

func (o Order) GetInfo(ctx context.Context, req *order.GetInfoRequest) (*order.GetInfoResponse, error) {
	mresp, err := o.service.GetOrderInfo(ctx, model.OrderID(req.OrderID))
	if err != nil {
		return nil, mapError(ctx, err)
	}

	resp := &order.GetInfoResponse{
		Status: order.OrderStatus(mresp.Status),
		User:   int64(mresp.UserID),
		Items:  make([]*order.Item, 0, len(mresp.Items)),
	}
	for _, item := range mresp.Items {
		resp.Items = append(resp.Items, &order.Item{
			Sku:   int32(item.SKU),
			Count: uint32(item.Count),
		})
	}

	return resp, nil
}

func (o Order) Pay(ctx context.Context, req *order.PayRequest) (*order.PayResponse, error) {
	if err := o.service.PayOrder(ctx, model.OrderID(req.OrderID)); err != nil {
		return nil, mapError(ctx, err)
	}
	return &order.PayResponse{}, nil
}

func (o Order) Cancel(ctx context.Context, req *order.CancelRequest) (*order.CancelResponse, error) {
	if err := o.service.CancelOrder(ctx, model.OrderID(req.OrderID)); err != nil {
		return nil, mapError(ctx, err)
	}
	return &order.CancelResponse{}, nil
}
