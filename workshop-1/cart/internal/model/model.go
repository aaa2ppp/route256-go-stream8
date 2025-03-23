package model

type (
	SKU     uint32
	UserID  int64
	OrderID int64
)

type Cart struct {
	UserID
	Items []CartItem
}

type CartItem struct {
	SKU
	Count uint16
}
