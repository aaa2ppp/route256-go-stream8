package model

type CreateOrderRequest struct {
	UserID UserID
	Items  []OrderItem
}

type CreateOrderResponse struct {
	OrderID OrderID
	Status  OrderStatus
}

type SetOrderStatusRequest = CreateOrderResponse
