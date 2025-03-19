package model

type (
	OrderID int64
	UserID  int64
	SKU     int64
)

type Order struct {
	OrderID
	UserID
	Status OrderStatus
	Items  []OrderItem
}

type OrderItem struct {
	SKU
	Count int
}

type Stock struct {
	SKU
	Count    int64
	Reserved int64
}
