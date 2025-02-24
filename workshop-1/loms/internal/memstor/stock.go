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
	count    uint64
	reserved uint64
}

func (si stockItem) available() uint64 {
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
			Count: rand.Uint64N(maxProductCount),
		})
	}

	stock := NewStock()
	stock.Init(items)
	return stock
}

type StockItem struct {
	SKU   model.SKU
	Count uint64
}

func (s *Stock) Init(items []StockItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = make(map[model.SKU]stockItem)
	for _, item := range items {
		s.items[item.SKU] = stockItem{count: item.Count}
	}
}

// GetInfo implements service.StockStorage.
func (s *Stock) GetInfo(_ context.Context, sku model.SKU) (count uint64, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[sku]
	if !exists {
		return 0, model.ErrNotFound
	}
	return item.available(), nil
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
		if stockItems[i].available() < uint64(orderItems[i].Count) {
			return fmt.Errorf("insufficient available SKU=%v", orderItems[i].SKU)
		}
		stockItems[i].reserved += uint64(orderItems[i].Count)
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
		if stockItems[i].reserved < uint64(orderItems[i].Count) {
			return fmt.Errorf("insufficient reserved SKU=%v", orderItems[i].SKU)
		}
		stockItems[i].reserved -= uint64(orderItems[i].Count)
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
		if stockItems[i].reserved < uint64(orderItems[i].Count) {
			return fmt.Errorf("insufficient reserved SKU=%v", orderItems[i].SKU)
		}
		n := uint64(orderItems[i].Count)
		stockItems[i].count -= n
		stockItems[i].reserved -= n
	}

	for i := range stockItems {
		s.items[orderItems[i].SKU] = stockItems[i]
	}

	return nil
}

var _ service.StockStorage = &Stock{}
