package client

import (
	"context"
	"net/http"
	"route256/cart/internal/config"
	"route256/cart/internal/model"
	"route256/cart/internal/service"
	"route256/cart/pkg/http/middleware"
)

type Product struct {
	client
	cfg *config.HTTPProductClient
}

func NewProduct(cfg *config.HTTPProductClient) Product {
	return Product{
		client: newClient(cfg.BaseURL, cfg.RequestTimeout),
		cfg:    cfg,
	}
}

type productGetInfoRequest struct {
	Token string    `json:"token"`
	SKU   model.SKU `json:"sku"`
}

type productGetInfoResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (c Product) GetInfo(ctx context.Context, sku model.SKU) (resp model.GetProductResponse, _ error) {
	log := getLogger(ctx, "Product.GetInfo")

	creq := productGetInfoRequest{
		Token: middleware.GetAuthTokenFromContext(ctx),
		SKU:   sku,
	}

	var cresp productGetInfoResponse
	status, err := c.doRequest(ctx, c.cfg.GetEndpoint, &creq, &cresp)
	if err != nil {
		log.Error("can't do request", "error", err)
		return resp, model.ErrInternalError
	}

	switch status {
	case http.StatusOK:
		return model.GetProductResponse{
			Name:  cresp.Name,
			Price: cresp.Price,
		}, nil
	case http.StatusUnauthorized:
		return resp, model.ErrUnauthorized
	case http.StatusNotFound:
		return resp, model.ErrNotFound
	}

	log.Error("unknown response status", "status", status)
	return resp, model.ErrInternalError
}

var _ service.ProductStorage = &Product{}
