package memstor

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
	"route256/loms/internal/service"
	"sync"
)

type Stock struct {
	mu    sync.RWMutex
	items map[model.SKU]model.Stock
}

func NewStock() *Stock {
	return &Stock{
		items: map[model.SKU]model.Stock{},
	}
}

func (s *Stock) Init(items []model.Stock) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = make(map[model.SKU]model.Stock, len(items))
	for i := range items {
		s.items[items[i].SKU] = items[i]
	}
}

// GetBySKU implements service.StockStorage.
func (s *Stock) GetBySKU(_ context.Context, sku model.SKU) (_ model.Stock, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[sku]
	if !exists {
		return model.Stock{}, model.ErrNotFound
	}
	return item, nil
}

// Reserve implements service.StockStorage.
func (s *Stock) Reserve(_ context.Context, orderItems []model.OrderItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stockItems := make([]model.Stock, 0, len(orderItems))
	for _, orderItem := range orderItems {
		stockItems = append(stockItems, s.items[orderItem.SKU])
	}

	for i := range stockItems {
		count := uint64(orderItems[i].Count)
		if stockItems[i].Available < count {
			return fmt.Errorf("insufficient available SKU=%v", orderItems[i].SKU)
		}
		stockItems[i].Available -= count
		stockItems[i].Reserved += count
	}

	for i := range stockItems {
		s.items[orderItems[i].SKU] = stockItems[i]
	}

	return nil
}

// ReserveCancel implements service.StockStorage.
func (s *Stock) ReserveCancel(_ context.Context, orderItems []model.OrderItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stockItems := make([]model.Stock, 0, len(orderItems))
	for _, orderItem := range orderItems {
		stockItems = append(stockItems, s.items[orderItem.SKU])
	}

	for i := range stockItems {
		count := uint64(orderItems[i].Count)
		if stockItems[i].Reserved < count {
			return fmt.Errorf("insufficient reserved SKU=%v", orderItems[i].SKU)
		}
		stockItems[i].Available += count
		stockItems[i].Reserved -= count
	}

	for i := range stockItems {
		s.items[orderItems[i].SKU] = stockItems[i]
	}

	return nil
}

// ReserveRemove implements service.StockStorage.
func (s *Stock) ReserveRemove(_ context.Context, orderItems []model.OrderItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stockItems := make([]model.Stock, 0, len(orderItems))
	for _, orderItem := range orderItems {
		stockItems = append(stockItems, s.items[orderItem.SKU])
	}

	for i := range stockItems {
		count := uint64(orderItems[i].Count)
		if stockItems[i].Reserved < count {
			return fmt.Errorf("insufficient reserved SKU=%v", orderItems[i].SKU)
		}
		stockItems[i].Reserved -= count
	}

	for i := range stockItems {
		s.items[orderItems[i].SKU] = stockItems[i]
	}

	return nil
}

var _ service.StockStorage = &Stock{}
