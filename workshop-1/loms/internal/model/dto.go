package model

type CreateOrderRequest struct {
	UserID UserID
	Items  []OrderItem
}

type SetOrderStatusRequest struct {
	OrderID OrderID
	Status  OrderStatus
}
