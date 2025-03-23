package model

type AddCartItemRequest struct {
	UserID
	SKU
	Count uint16
}

type DeleteCartItemRequest struct {
	UserID
	SKU
}

type GetProductResponse struct {
	Name  string
	Price uint32
}

type CartListResponse struct {
	Items      []CartListItem
	TotalPrice uint32
}

type CartListItem struct {
	SKU   SKU
	Count uint16
	Name  string
	Price uint32
}

type OrderCreateRequest struct {
	UserID UserID
	Items  []OrderItem
}

type OrderItem = CartItem
