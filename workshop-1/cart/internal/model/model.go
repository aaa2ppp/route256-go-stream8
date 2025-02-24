package model

type (
	SKU     int32
	UserID  int64
	OrderID int64
)

type Cart struct {
	UserID UserID
	Items  []CartItem
}

type CartItem struct {
	SKU   SKU
	Count uint16
}
