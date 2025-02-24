package handler

import (
	"context"
	"errors"
	"net/http"
	"route256/loms/internal/model"
	"strconv"
)

type orderItem struct {
	SKU   model.SKU `json:"sku"`
	Count uint16    `json:"count"`
}

type createOrderRequest struct {
	User  model.UserID `json:"user"`
	Items []orderItem  `json:"items"`
}

func (req createOrderRequest) Validate() error {
	var errs []error
	if req.User <= 0 {
		errs = append(errs, errors.New("user must be > 0"))
	}
	if len(req.Items) == 0 {
		errs = append(errs, errors.New("items cannot be empty"))
	}
	for i, item := range req.Items {
		if item.SKU <= 0 {
			errs = append(errs, errors.New("item["+strconv.Itoa(i)+"]: sku must be > 0"))
		}
		if item.Count <= 0 {
			errs = append(errs, errors.New("item["+strconv.Itoa(i)+"]: count must be > 0"))
		}
	}
	return errors.Join(errs...)
}

type createOreserResponce struct {
	OrderID model.OrderID `json:"orderID"`
}

func CreateOrder(crateFunc func(ctx context.Context, req model.CreateOrderRequest) (resp model.CreateOrderResponse, err error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "CreateOrder")

		var req createOrderRequest
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		mreq := model.CreateOrderRequest{
			UserID: req.User,
			Items:  make([]model.OrderItem, 0, len(req.Items)),
		}
		for _, item := range req.Items {
			mreq.Items = append(mreq.Items, model.OrderItem{
				SKU:   item.SKU,
				Count: item.Count,
			})
		}

		mresp, err := crateFunc(x.ctx(), mreq)
		if err != nil {
			x.writeError(err)
			return
		}

		x.writeResponse(200, createOreserResponce{mresp.OrderID}) // or 201 ?
	}
}

type getOrderRequest struct {
	OrderID model.OrderID `json:"orderID"`
}

type getOrderInfoRequest = getOrderRequest

func (req getOrderInfoRequest) Validate() error {
	if req.OrderID <= 0 {
		return errors.New("orderID must be > 0")
	}
	return nil
}

type getOrderInfoResponse struct {
	Status string       `json:"status"`
	User   model.UserID `json:"user"`
	Items  []orderItem  `json:"items"`
}

func GetOrderInfo(getFunc func(ctx context.Context, orderID model.OrderID) (resp model.Order, err error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "GetOrederInfo")

		var req getOrderInfoRequest
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		mresp, err := getFunc(x.ctx(), req.OrderID)
		if err != nil {
			x.writeError(err)
			return
		}

		resp := getOrderInfoResponse{
			Status: mresp.Status.String(),
			User:   mresp.UserID,
			Items:  make([]orderItem, 0, len(mresp.Items)),
		}
		for _, item := range mresp.Items {
			resp.Items = append(resp.Items, orderItem{
				SKU:   item.SKU,
				Count: item.Count,
			})
		}

		x.writeResponse(200, resp)
	}
}

type payOrderRequest = getOrderRequest

func PayOrder(payFunc func(ctx context.Context, orderID model.OrderID) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "PayOrder")

		var req payOrderRequest
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		if err := payFunc(x.ctx(), req.OrderID); err != nil {
			x.writeError(err)
			return
		}

		x.writeResponse(200, struct{}{})
	}
}

type cancelOrderRequest = getOrderRequest

func CancelOrder(cancelFunc func(ctx context.Context, orderID model.OrderID) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "CancelOrder")

		var req cancelOrderRequest
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		if err := cancelFunc(x.ctx(), req.OrderID); err != nil {
			x.writeError(err)
			return
		}

		x.writeResponse(200, struct{}{})
	}
}

type getStockInfoRequest struct {
	SKU model.SKU `json:"sku"`
}

func (req getStockInfoRequest) Validate() error {
	if req.SKU <= 0 {
		return errors.New("sku must be > 0")
	}
	return nil
}

type getStockInfoResponse struct {
	Count uint64 `json:"count"`
}

func GetStockInfo(getFunc func(ctx context.Context, sku model.SKU) (count uint64, err error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := newHelper(w, r, "GetStockInfo")

		var req getStockInfoRequest
		if !x.decodeBodyAndValidateRequest(&req) {
			return
		}

		count, err := getFunc(x.ctx(), req.SKU)
		if err != nil {
			x.writeError(err)
			return
		}

		x.writeResponse(200, getStockInfoResponse{count})
	}
}
