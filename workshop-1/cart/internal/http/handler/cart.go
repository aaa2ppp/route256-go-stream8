package handler

import (
	"context"
	"errors"
	"net/http"
	"route256/cart/internal/model"
)

type cartAddItemRequest struct {
	User  model.UserID `json:"user"`
	SKU   model.SKU    `json:"sku"`
	Count uint16       `json:"count"`
}

func (req cartAddItemRequest) Validate() error {
	var errs []error
	if req.User <= 0 {
		errs = append(errs, errors.New("user must be > 0"))
	}
	if req.SKU <= 0 {
		errs = append(errs, errors.New("sku must be > 0"))
	}
	return errors.Join(errs...)
}

func CartAddItem(addFunc func(ctx context.Context, req model.AddCartItemRequest) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "CartAddItem")

		if !x.checkPOSTMethod() {
			return
		}

		var req cartAddItemRequest
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		if err := addFunc(x.ctx(), model.AddCartItemRequest{
			UserID: req.User,
			SKU:    req.SKU,
			Count:  req.Count,
		}); err != nil {
			x.writeError(err)
			return
		}

		x.writeResponse(http.StatusOK, struct{}{})
	}
}

type cartDeleteItemRequest struct {
	User model.UserID `json:"user"`
	SKU  model.SKU    `json:"sku"`
}

func (req cartDeleteItemRequest) Validate() error {
	var errs []error
	if req.User <= 0 {
		errs = append(errs, errors.New("user must be > 0"))
	}
	if req.SKU <= 0 {
		errs = append(errs, errors.New("sku must be > 0"))
	}
	return errors.Join(errs...)
}

func CartDeleteItem(deleteFunc func(ctx context.Context, req model.DeleteCartItemRequest) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "CartDeleteItem")

		if !x.checkPOSTMethod() {
			return
		}

		var req cartDeleteItemRequest
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		if err := deleteFunc(x.ctx(), model.DeleteCartItemRequest{
			UserID: req.User,
			SKU:    req.SKU,
		}); err != nil {
			x.writeError(err)
			return
		}

		x.writeResponse(http.StatusOK, struct{}{})
	}
}

type userCartRequest struct {
	User model.UserID `json:"user"`
}

func (req userCartRequest) Validate() error {
	if req.User <= 0 {
		return errors.New("user must be > 0")
	}
	return nil
}

type cartList = userCartRequest

type cartListItem struct {
	SKU   model.SKU `json:"sku"`
	Count uint16    `json:"count"`
	Name  string    `json:"name"`
	Price uint32    `json:"price"`
}

type cartListResponse struct {
	Items      []cartListItem `json:"items"`
	TotalPrice uint32         `json:"totalPrice"`
}

func CartList(listFunc func(ctx context.Context, userID model.UserID) (model.CartListResponse, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "CartList")

		if !x.checkPOSTMethod() {
			return
		}

		var req cartList
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		cart, err := listFunc(x.ctx(), req.User)
		if err != nil {
			x.writeError(err)
			return
		}

		resp := cartListResponse{
			Items:      make([]cartListItem, 0, len(cart.Items)),
			TotalPrice: 0,
		}
		for _, item := range cart.Items {
			resp.Items = append(resp.Items, cartListItem{
				SKU:   item.SKU,
				Count: item.Count,
				Name:  item.Name,
				Price: item.Price,
			})
		}
		resp.TotalPrice = cart.TotalPrice

		x.writeResponse(http.StatusOK, resp)
	}
}

type cartClearRequest = userCartRequest

func CartClear(clearFunc func(ctx context.Context, userID model.UserID) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "CartClear")

		if !x.checkPOSTMethod() {
			return
		}

		var req cartClearRequest
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		if err := clearFunc(x.ctx(), req.User); err != nil {
			x.writeError(err)
			return
		}

		x.writeResponse(http.StatusOK, struct{}{})
	}
}

type cartCheckoutRequest = userCartRequest

type cartCheckoutResponse struct {
	OrderID model.OrderID `json:"orderID"`
}

func CartCheckout(checkoutFunc func(ctx context.Context, userID model.UserID) (model.OrderID, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "CartCheckout")

		if !x.checkPOSTMethod() {
			return
		}

		var req cartCheckoutRequest
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		orderID, err := checkoutFunc(x.ctx(), req.User)
		if err != nil {
			x.writeError(err)
			return
		}

		x.writeResponse(http.StatusOK, cartCheckoutResponse{
			OrderID: orderID,
		})
	}
}
