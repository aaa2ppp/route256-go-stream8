package memstor

import (
	"context"
	"route256/cart/internal/model"
	"route256/cart/internal/service"
	"slices"
	"sync"
)

type Cart struct {
	mu    sync.RWMutex
	carts map[model.UserID][]model.CartItem
}

func NewCart() *Cart {
	return &Cart{
		carts: map[model.UserID][]model.CartItem{},
	}
}

func (p *Cart) Add(_ context.Context, req model.AddCartItemRequest) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cart := p.carts[req.UserID]

reqItemsLoop:
	for _, reqItem := range req.Items {
		for i := range cart {
			item := &cart[i]
			if item.SKU == reqItem.SKU {
				item.Count += reqItem.Count
				continue reqItemsLoop
			}
		}
		cart = append(cart, model.CartItem{SKU: reqItem.SKU, Count: reqItem.Count})
	}

	p.carts[req.UserID] = cart
	return nil
}

func (p *Cart) Delete(_ context.Context, req model.DeleteCartItemRequest) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cart := p.carts[req.UserID]
	for i := range cart {
		item := &cart[i]
		if item.SKU == req.SKU {
			n := len(cart)
			cart[i] = cart[n-1]
			cart = cart[:n-1]
			p.carts[req.UserID] = cart
			return nil
		}
	}

	return model.ErrNotFound
}

func (p *Cart) List(_ context.Context, userID model.UserID) ([]model.CartItem, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	cart := p.carts[userID]
	if cart == nil {
		return nil, model.ErrNotFound
	}

	return slices.Clone(cart), nil
}

func (p *Cart) Clear(_ context.Context, userID model.UserID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.carts, userID)
	return nil
}

var _ service.CartStorage = &Cart{}
