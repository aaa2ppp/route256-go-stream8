package model

type (
	OrderID int64
	UserID  int64
	SKU     uint32
)

type Order struct {
	OrderID
	UserID
	Status OrderStatus
	Items  []OrderItem
}

type OrderItem struct {
	SKU
	Count uint16
}

type Stock struct {
	SKU
	Available uint64
	Reserved  uint64
}
