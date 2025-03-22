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

	items := p.carts[req.UserID]

	for i, item := range items {
		if item.SKU == req.SKU {
			items[i].Count = req.Count
			return nil
		}
	}

	items = append(items, model.CartItem{SKU: req.SKU, Count: req.Count})
	p.carts[req.UserID] = items

	return nil
}

func (p *Cart) Delete(_ context.Context, req model.DeleteCartItemRequest) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	items := p.carts[req.UserID]

	for i := range items {
		if items[i].SKU == req.SKU {
			p.carts[req.UserID] = append(items[:i], items[i+1:]...)
			return nil
		}
	}

	return model.ErrNotFound
}

func (p *Cart) List(_ context.Context, userID model.UserID) ([]model.CartItem, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	items := p.carts[userID]

	if items == nil {
		return nil, model.ErrNotFound
	}

	return slices.Clone(items), nil
}

func (p *Cart) Clear(_ context.Context, userID model.UserID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.carts, userID)
	return nil
}

var _ service.CartStorage = &Cart{}
