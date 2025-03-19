package memstor

import (
	"context"
	"fmt"
	"math/rand/v2"
	"route256/loms/internal/model"
	"route256/loms/internal/service"
	"sync"
)

type stockItem struct {
	count    int64
	reserved int64
}

func (si stockItem) available() int64 {
	return si.count - si.reserved
}

type Stock struct {
	mu    sync.RWMutex
	items map[model.SKU]stockItem
}

func NewStock() *Stock {
	return &Stock{
		items: map[model.SKU]stockItem{},
	}
}

func NewRandomStock(producs []model.SKU) *Stock {
	const maxProductCount = 10

	items := make([]StockItem, 0, len(producs))
	for _, sku := range producs {
		items = append(items, StockItem{
			SKU:   sku,
			Count: rand.Int64N(maxProductCount),
		})
	}

	stock := NewStock()
	stock.Init(items)
	return stock
}

type StockItem struct {
	SKU   model.SKU
	Count int64
}

func (s *Stock) Init(items []StockItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = make(map[model.SKU]stockItem)
	for _, item := range items {
		s.items[item.SKU] = stockItem{count: item.Count}
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
	return model.Stock{
		SKU:      sku,
		Count:    int64(item.count),
		Reserved: int64(item.reserved),
	}, nil
}

// Reserve implements service.StockStorage.
func (s *Stock) Reserve(_ context.Context, orderItems []model.OrderItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stockItems := make([]stockItem, 0, len(orderItems))
	for _, orderItem := range orderItems {
		stockItems = append(stockItems, s.items[orderItem.SKU])
	}

	for i := range stockItems {
		n := int64(orderItems[i].Count)
		if stockItems[i].available() < n {
			return fmt.Errorf("insufficient available SKU=%v", orderItems[i].SKU)
		}
		stockItems[i].reserved += n
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

	stockItems := make([]stockItem, 0, len(orderItems))
	for _, orderItem := range orderItems {
		stockItems = append(stockItems, s.items[orderItem.SKU])
	}

	for i := range stockItems {
		n := int64(orderItems[i].Count)
		if stockItems[i].reserved < n {
			return fmt.Errorf("insufficient reserved SKU=%v", orderItems[i].SKU)
		}
		stockItems[i].reserved -= n
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

	stockItems := make([]stockItem, 0, len(orderItems))
	for _, orderItem := range orderItems {
		stockItems = append(stockItems, s.items[orderItem.SKU])
	}

	for i := range stockItems {
		n := int64(orderItems[i].Count)
		if stockItems[i].reserved < n {
			return fmt.Errorf("insufficient count SKU=%v", orderItems[i].SKU)
		}
		if stockItems[i].reserved < n {
			return fmt.Errorf("insufficient reserved SKU=%v", orderItems[i].SKU)
		}
		stockItems[i].count -= n
		stockItems[i].reserved -= n
	}

	for i := range stockItems {
		s.items[orderItems[i].SKU] = stockItems[i]
	}

	return nil
}

var _ service.StockStorage = &Stock{}
