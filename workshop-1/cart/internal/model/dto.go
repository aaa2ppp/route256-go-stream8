package model

type AddCartItemRequest struct {
	Token  string
	UserID UserID
	Items  []CartItem
}

type DeleteCartItemRequest struct {
	UserID UserID
	SKU    SKU
}

type GetProductRequest struct {
	Token string
	SKU   SKU
}

type GetProductResponse struct {
	Name  string
	Price uint32
}

type CartListItem struct {
	SKU   SKU
	Count uint16
	Name  string
	Price uint32
}

type CartListRequest struct {
	UserID UserID
	Token  string
}

type CartListResponse struct {
	Items      []CartListItem
	TotalPrice uint32
}

type OrderItem = CartItem

type OrderCreateRequest struct {
	UserID UserID
	Items  []OrderItem
}
