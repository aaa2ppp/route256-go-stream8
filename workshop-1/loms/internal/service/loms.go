package service

import (
	"context"
	"route256/loms/internal/model"
)

type OrderStorage interface {
	Create(ctx context.Context, req model.CreateOrderRequest) (model.OrderID, error)
	SetStatus(ctx context.Context, req model.SetOrderStatusRequest) error
	GetByID(ctx context.Context, orderID model.OrderID) (model.Order, error)
}

type StockStorage interface {
	Reserve(ctx context.Context, items []model.OrderItem) error
	ReserveRemove(ctx context.Context, items []model.OrderItem) error
	ReserveCancel(ctx context.Context, items []model.OrderItem) error
	GetBySKU(ctx context.Context, sku model.SKU) (stock model.Stock, err error)
}

type LOMS struct {
	order OrderStorage
	stock StockStorage
}

func NewLOMS(order OrderStorage, stock StockStorage) *LOMS {
	return &LOMS{
		order: order,
		stock: stock,
	}
}

func (p *LOMS) CreateOrder(ctx context.Context, req model.CreateOrderRequest) (orderID model.OrderID, err error) {
	log := getLogger(ctx, "LOMS.CreateOrder")

	orderID, err = p.order.Create(ctx, req)
	if err != nil {
		return orderID, err
	}

	var status model.OrderStatus
	if err := p.stock.Reserve(ctx, req.Items); err != nil {
		status = model.OrderStatusFailed
	} else {
		status = model.OrderStatusAwaitingPayment
	}

	if err := p.order.SetStatus(ctx, model.SetOrderStatusRequest{
		OrderID: orderID,
		Status:  status,
	}); err != nil {
		log.Error("can't set order status", "error", err)
		return 0, model.ErrInternalError
	}

	if status == model.OrderStatusAwaitingPayment {
		return orderID, nil
	} else {
		return 0, model.ErrPreconditionFailed
	}
}

func (p *LOMS) GetOrderInfo(ctx context.Context, orderID model.OrderID) (resp model.Order, err error) {
	return p.order.GetByID(ctx, orderID)
}

func (p *LOMS) PayOrder(ctx context.Context, orderID model.OrderID) error {
	log := getLogger(ctx, "LOMS.PayOrder")

	orderInfo, err := p.order.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if orderInfo.Status != model.OrderStatusAwaitingPayment {
		log.Debug("order status must be 'awaiting payment'", "orderStatus", orderInfo.Status)
		return model.ErrPreconditionFailed
	}

	if err := p.stock.ReserveRemove(ctx, orderInfo.Items); err != nil {
		log.Error("can't reserve remove", "error", err)
		return model.ErrPreconditionFailed
	}

	if err := p.order.SetStatus(ctx, model.SetOrderStatusRequest{
		OrderID: orderID,
		Status:  model.OrderStatusPayed,
	}); err != nil {
		log.Error("can't set order status", "error", err)
		return model.ErrInternalError
	}

	return nil
}

func (p *LOMS) CancelOrder(ctx context.Context, orderID model.OrderID) error {
	log := getLogger(ctx, "LOMS.CancelOrder")

	orderInfo, err := p.order.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if orderInfo.Status != model.OrderStatusAwaitingPayment {
		log.Debug("order status must be 'awaiting payment'", "orderStatus", orderInfo.Status)
		return model.ErrPreconditionFailed
	}

	if err := p.stock.ReserveCancel(ctx, orderInfo.Items); err != nil {
		log.Error("can't reserve cancel", "error", err)
		return model.ErrPreconditionFailed
	}

	if err := p.order.SetStatus(ctx, model.SetOrderStatusRequest{
		OrderID: orderID,
		Status:  model.OrderStatusCancelled,
	}); err != nil {
		log.Error("can't set order status", "error", err)
		return model.ErrInternalError
	}

	return nil
}

func (p *LOMS) GetStockInfo(ctx context.Context, sku model.SKU) (count uint64, err error) {
	stock, err := p.stock.GetBySKU(ctx, sku)
	if err != nil {
		return 0, err
	}
	return uint64(stock.Available), nil
}
