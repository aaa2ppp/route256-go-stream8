package model

type (
	UserID  int64
	SKU     int32
	OrderID int64
)

type OrderItem struct {
	SKU   SKU
	Count uint16
}

type Order struct {
	Status OrderStatus
	UserID
	Items []OrderItem
}
